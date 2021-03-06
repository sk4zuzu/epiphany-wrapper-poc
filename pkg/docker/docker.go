/*
 * Copyright © 2020 Mateusz Kyc
 */

package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Image struct {
	Name string
}

func (i *Image) Pull() (string, error) { //TODO consider passing log file directly here to write to it on the fly
	ctx, cli, err := clientAndContext()
	if err != nil {
		return "", err
	}
	reader, err := cli.ImagePull(ctx, i.Name, types.ImagePullOptions{}) //TODO format output
	logR, logW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	done := make(chan bool)
	defer close(done)

	var result string

	go func() {
		_, _ = io.Copy(os.Stdout, stdoutR)
		done <- true
	}()

	go func() {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, logR)
		result = buf.String()
		done <- true
	}()

	go func() {
		defer logW.Close()
		defer stdoutW.Close()

		// build the MultiWriter for all the pipes
		mw := io.MultiWriter(logW, stdoutW)

		// copy the data into the MultiWriter
		_, _ = io.Copy(mw, reader)
	}()

	for c := 0; c < 2; c++ {
		<-done
	}

	reader.Close()
	if err != nil {
		return result, err
	}
	return result, nil
}

type Job struct {
	Image                string
	Command              string
	Args                 []string
	WorkDirectory        string
	Mounts               []string
	MountPath            string
	EnvironmentVariables map[string]string
}

func (j Job) Run() error {
	return run(j)
}

func run(job Job) error {
	ctx, cli, err := clientAndContext()
	if err != nil {
		return err
	}
	var envs []string
	for k, v := range job.EnvironmentVariables {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	commandAndArgs := append([]string{job.Command}, job.Args...)
	var mounts []mount.Mount
	for _, m := range job.Mounts {
		mounts = append(
			mounts,
			mount.Mount{
				Type:   mount.TypeBind,
				Source: path.Join(job.MountPath, m),
				Target: m,
			})
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:      job.Image,
			Cmd:        commandAndArgs,
			WorkingDir: job.WorkDirectory,
			Env:        envs,
			Tty:        false,
		}, &container.HostConfig{
			Mounts: mounts,
		},
		nil,
		"",
	)
	if err != nil {
		return err
	}
	defer cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}) //TODO warn if error while removing

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		return err
	}

	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out) //TODO write logs to file as well

	return nil
}

func clientAndContext() (context.Context, *client.Client, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, nil, err
	}
	return ctx, cli, nil
}

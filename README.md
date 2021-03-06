# epiphany-wrapper-poc
PoC application to wrap containerised modules of epiphany

## example run

### general 

#### e help 

```shell
> e --help                        
E wrapper allows to interact with epiphany

Usage:
  e [command]

Available Commands:
  components   Allows to inspect and install available components
  environments Allows various interactions with environments
  help         Help about any command

Flags:
      --config string   config file (default is /Users/mateusz/.e/config.yaml)
  -h, --help            help for e

Use "e [command] --help" for more information about a command.

```

### components sub-command

#### e components help

```shell
> e components --help
This command provides way to:
 - list available components, 
 - install new component to environment
 - get information about component

Information about available components are taken from https://github.com/mkyc/epiphany-wrapper-poc-repo/blob/master/v1.yaml

Usage:
  e components [flags]
  e components [command]

Available Commands:
  info        Displays information about component
  install     Installs component into currently used environment
  list        Lists all existing components in repository

Flags:
  -h, --help   help for components

Global Flags:
      --config string   config file (default is /Users/mateusz/.e/config.yaml)

Use "e components [command] --help" for more information about a command.

```

#### e components list 

```shell
> e components list  
Using config file: /Users/mateusz/.e/config.yaml
list called
c1 0.1.0
luuk-c1 0.0.1
```

#### e components info

```shell
> e components info c1
Using config file: /Users/mateusz/.e/config.yaml
components info called
Component:
 Name: c1
 Type: docker
 Version: 0.1.0
 Image: docker.io/hashicorp/terraform:0.12.28
 Commands:
  Name: init
  Description: initializes terraform in local directory
  Name: apply
  Description: applies terraform in local directory
```

#### e components install

Use you created environment with `e environment new e1`

```shell
> e components install c1
Using config file: /Users/mateusz/.e/config.yaml
install called
{"status":"Pulling from hashicorp/terraform","id":"0.12.28"}
{"status":"Digest: sha256:19a84a76564c9bea081b405f458a51107cf2abbafba34af2f02774e72d551ad1"}
{"status":"Status: Image is up to date for hashicorp/terraform:0.12.28"}
Installed component c1 0.1.0 to environment e1
```

### environments sub-command

#### e environments help

```shell
> e environments --help
TODO

Usage:
  e environments [flags]
  e environments [command]

Available Commands:
  info        Displays information about currently selected environment
  new         Creates new environment
  run         Runs installed component command in environment
  use         Allows to select environment to be used

Flags:
  -h, --help   help for environments

Global Flags:
      --config string   config file (default is /Users/mateusz/.e/config.yaml)

Use "e environments [command] --help" for more information about a command.
```

#### e environments new

```shell
> e environments new e1
Using config file: /Users/mateusz/.e/config.yaml
new called
new created environment is: e1
```

#### e environments info

Here used after command `e components install c1`

```shell
> e environments info    
Using config file: /Users/mateusz/.e/config.yaml
enviroments info called
Environment info:
 Name: e1
 UUID: ade1b8ad-3723-4f85-b51a-3cffa057b2c8
  Installed Component:
   Name: c1
   Type: docker
   Version: 0.1.0
   Image: docker.io/hashicorp/terraform:0.12.28
    Command:
     Name init
     Description initializes terraform in local directory
    Command:
     Name apply
     Description applies terraform in local directory
```

#### e environments use

```shell 
> e environments use
Using config file: /Users/mateusz/.e/config.yaml
use called
✔ e1 (ade1b8ad-3723-4f85-b51a-3cffa057b2c8, current)
Choosen environment UUID is: ade1b8ad-3723-4f85-b51a-3cffa057b2c8
```

#### e environments run

```shell
> e environments run c1 init
Using config file: /Users/mateusz/.e/config.yaml
run called
2020/07/28 15:39:17 [WARN] Log levels other than TRACE are currently unreliable, and are supported only for backward compatibility.
  Use TF_LOG=TRACE to see Terraform's internal logs.
  ----
Terraform initialized in an empty directory!

The directory has no Terraform configuration files. You may begin working
with Terraform immediately by creating Terraform configuration files.
running command completed!
```

## configuration directory structure

After all command executed in previous section directory structure looks in similar way to: 

```shell
> pwd
/Users/mateusz/.e
> tree -a
.
├── config.yaml
├── environments
│   └── ade1b8ad-3723-4f85-b51a-3cffa057b2c8
│       ├── c1
│       │   └── 0.1.0
│       │       ├── mounts
│       │       │   └── terraform
│       │       └── runs
│       │           └── 20200728-173415.915CEST.log
│       └── config.yaml
└── v1.yaml

7 directories, 4 files
```

Main config file contains: 

```yaml
> cat config.yaml 
version: v1
kind: Config
current-environment: ade1b8ad-3723-4f85-b51a-3cffa057b2c8
```

Used environment config file contains: 

```yaml
> cat environments/ade1b8ad-3723-4f85-b51a-3cffa057b2c8/config.yaml 
name: e1
uuid: ade1b8ad-3723-4f85-b51a-3cffa057b2c8
installed:
- environment_ref: ade1b8ad-3723-4f85-b51a-3cffa057b2c8
  name: c1
  type: docker
  version: 0.1.0
  image: docker.io/hashicorp/terraform:0.12.28
  workdir: /terraform
  mounts:
  - /terraform
  commands:
  - name: init
    description: initializes terraform in local directory
    command: init
    envs:
      TF_LOG: WARN
    args: []
  - name: apply
    description: applies terraform in local directory
    command: apply
    envs:
      TF_LOG: DEBUG
    args:
    - -auto-approve
```

## TODO

There is a lot TODO's in a code which should be fixed
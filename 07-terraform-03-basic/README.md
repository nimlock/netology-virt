# Домашнее задание к занятию "7.3. Основы и принцип работы Терраформ"

## Модуль 7. Облачная инфраструктура. Terraform

### Студент: Иван Жиляев

## Задача 1. Создадим бэкэнд в S3 (необязательно, но крайне желательно).

>Если в рамках предыдущего задания у вас уже есть аккаунт AWS, то давайте продолжим знакомство со взаимодействием
>терраформа и aws. 
>
>1. Создайте s3 бакет, iam роль и пользователя от которого будет работать терраформ. Можно создать отдельного пользователя,
>а можно использовать созданного в рамках предыдущего задания, просто добавьте ему необходимы права, как описано 
>[здесь](https://www.terraform.io/docs/backends/types/s3.html).
>1. Зарегистрируйте бэкэнд в терраформ проекте как описано по ссылке выше. 

1. Буду использовать пользователя `terraform`, созданного в [предыдущем задании](../07-terraform-02-syntax/README.md).  
S3 bucket создал через web-консоль AWS, в этом хранилище создал папку `terraform-remote-config/` и выставил требуемые права для bucket следующей политикой:

    ```
    {
        "Version": "2012-10-17",
        "Id": "Policy1607622331907",
        "Statement": [
            {
                "Sid": "Stmt1607622315031",
                "Effect": "Allow",
                "Principal": {
                    "AWS": "arn:aws:iam::536867559799:user/terraform"
                },
                "Action": [
                    "s3:GetObject",
                    "s3:PutObject"
                ],
                "Resource": "arn:aws:s3:::ivan-07-terraform-bucket/terraform-remote-config"
            },
            {
                "Sid": "Stmt1607622713110",
                "Effect": "Allow",
                "Principal": {
                    "AWS": "arn:aws:iam::536867559799:user/terraform"
                },
                "Action": "s3:ListBucket",
                "Resource": "arn:aws:s3:::ivan-07-terraform-bucket"
            }
        ]
    }
    ```

1. Бэкенд в терраформ зарегистрирован благодаря добавлению в [main.tf](https://github.com/nimlock/netology-terraform_repository/blob/e126472559ee4915286416442eb02c51e8fad18d/main.tf) такого блока:

    ```
    terraform {
    backend "s3" {
        bucket = "ivan-07-terraform-bucket"
        key    = "terraform-remote-config/terraform.tfstate"
        region = "us-east-2"
        dynamodb_table = "terraform-state-locking"
    }
    }
    ```

## Задача 2. Инициализируем проект и создаем воркспейсы. 

>1. Выполните `terraform init`:
>    * если был создан бэкэнд в S3, то терраформ создат файл стейтов в S3 и запись в таблице 
>dynamodb.
>    * иначе будет создан локальный файл со стейтами.  
>1. Создайте два воркспейса `stage` и `prod`.
>1. В уже созданный `aws_instance` добавьте зависимость типа инстанса от вокспейса, что бы в разных ворскспейсах 
>использовались разные `instance_type`.
>1. Добавим `count`. Для `stage` должен создаться один экземпляр `ec2`, а для `prod` два. 
>1. Создайте рядом еще один `aws_instance`, но теперь определите их количество при помощи `for_each`, а не `count`.
>1. Что бы при изменении типа инстанса не возникло ситуации, когда не будет ни одного инстанса добавьте параметр
>жизненного цикла `create_before_destroy = true` в один из рессурсов `aws_instance`.
>1. При желании поэкспериментируйте с другими параметрами и рессурсами.
>
>В виде результата работы пришлите:
>* Вывод команды `terraform workspace list`.
>* Вывод команды `terraform plan` для воркспейса `prod`.

1. Для создания стейта именно в S3 необходимо выполнить `terraform apply`, ведь только при непосредственном применении изменений к инфраструктуре мы имеем право на фиксацию общего стейта.

1. Создадим воркспейсы:

    ```
    terraform workspace new stage
    terraform workspace new prod
    ```

1. Для реализации зависимости нашей конфигурации от воркспейсов необходимо добавить блок, задающий локальную переменную (словарь) в которой опишем желаемые зависимости:

    ```
    locals {
      dict_of_instance_types = {
          stage = "t2.micro"
          prod = "t2.small"
      }
    }
    ```

    Осталось заменить явное значение параметра создаваемого ресурса ссылкой на переменную:

    ```
    resource "aws_instance" "web" {
      ami           = data.aws_ami.ubuntu.id
      instance_type = locals.dict_of_instance_types[terraform.workspace]
    ...
    ```

1. Для добавления зависимости количества инстансов от окружения добавим ещё один блок с описанием ещё одного словаря:

    ```
    locals {
      dict_of_instance_count = {
          stage = 1
          prod = 2
      }
    }
    ```

    Теперь используем эту переменную в описании ресурса:

    ```
    resource "aws_instance" "web" {
    ami           = data.aws_ami.ubuntu.id
    instance_type = locals.dict_of_instance_types[terraform.workspace]
    count = locals.dict_of_instance_count[terraform.workspace]
    ...
    ```

1. В случае работы с `for_each` нам опять же потребуется список (любой итерируемый объект), зададим его по образу примера с лекции:

    ```
    locals {
      instances = {
          "t2.micro" = data.aws_ami.ubuntu.id
          "t2.small" = data.aws_ami.ubuntu.id
      }
    }
    ```

    В описании нового ресурса типа `aws_instance` применим итерацию по словарю, при этом параметры ресурса определяются из записей словаря:

    ```
    resource "aws_instance" "api" {
    for_each = local.instances

    ami           = each.value
    instance_type = each.key
    }
    ```

1. Для изменения жизненного цикла одного из ресурсов добавим в него такую директиву:

    ```
    ...
    lifecycle {
      create_before_destroy = true
    }
    ...
    ```

### Результаты выполнения задания

- получившаяся конфигурация находится в файле [main.tf](https://github.com/nimlock/netology-terraform_repository/blob/e126472559ee4915286416442eb02c51e8fad18d/main.tf)

- ```
  ivan@kubang:~/study/netology-virt/terraform_repository$ terraform validate
  Success! The configuration is valid.
  ```

- ```
  ivan@kubang:~/study/netology-virt/terraform_repository$ terraform workspace list
    default
  * prod
    stage
  ```

- ```
  ivan@kubang:~/study/netology-virt/terraform_repository$ terraform plan
  Acquiring state lock. This may take a few moments...
  
  An execution plan has been generated and is shown below.
  Resource actions are indicated with the following symbols:
    + create
  
  Terraform will perform the following actions:
  
    # aws_instance.api["t2.micro"] will be created
    + resource "aws_instance" "api" {
        + ami                          = "ami-0b289b3e97908e84e"
        + arn                          = (known after apply)
        + associate_public_ip_address  = (known after apply)
        + availability_zone            = (known after apply)
        + cpu_core_count               = (known after apply)
        + cpu_threads_per_core         = (known after apply)
        + get_password_data            = false
        + host_id                      = (known after apply)
        + id                           = (known after apply)
        + instance_state               = (known after apply)
        + instance_type                = "t2.micro"
        + ipv6_address_count           = (known after apply)
        + ipv6_addresses               = (known after apply)
        + key_name                     = (known after apply)
        + outpost_arn                  = (known after apply)
        + password_data                = (known after apply)
        + placement_group              = (known after apply)
        + primary_network_interface_id = (known after apply)
        + private_dns                  = (known after apply)
        + private_ip                   = (known after apply)
        + public_dns                   = (known after apply)
        + public_ip                    = (known after apply)
        + secondary_private_ips        = (known after apply)
        + security_groups              = (known after apply)
        + source_dest_check            = true
        + subnet_id                    = (known after apply)
        + tenancy                      = (known after apply)
        + volume_tags                  = (known after apply)
        + vpc_security_group_ids       = (known after apply)
  
        + ebs_block_device {
            + delete_on_termination = (known after apply)
            + device_name           = (known after apply)
            + encrypted             = (known after apply)
            + iops                  = (known after apply)
            + kms_key_id            = (known after apply)
            + snapshot_id           = (known after apply)
            + volume_id             = (known after apply)
            + volume_size           = (known after apply)
            + volume_type           = (known after apply)
          }
  
        + ephemeral_block_device {
            + device_name  = (known after apply)
            + no_device    = (known after apply)
            + virtual_name = (known after apply)
          }
  
        + metadata_options {
            + http_endpoint               = (known after apply)
            + http_put_response_hop_limit = (known after apply)
            + http_tokens                 = (known after apply)
          }
  
        + network_interface {
            + delete_on_termination = (known after apply)
            + device_index          = (known after apply)
            + network_interface_id  = (known after apply)
          }
  
        + root_block_device {
            + delete_on_termination = (known after apply)
            + device_name           = (known after apply)
            + encrypted             = (known after apply)
            + iops                  = (known after apply)
            + kms_key_id            = (known after apply)
            + volume_id             = (known after apply)
            + volume_size           = (known after apply)
            + volume_type           = (known after apply)
          }
      }
  
    # aws_instance.api["t2.small"] will be created
    + resource "aws_instance" "api" {
        + ami                          = "ami-0b289b3e97908e84e"
        + arn                          = (known after apply)
        + associate_public_ip_address  = (known after apply)
        + availability_zone            = (known after apply)
        + cpu_core_count               = (known after apply)
        + cpu_threads_per_core         = (known after apply)
        + get_password_data            = false
        + host_id                      = (known after apply)
        + id                           = (known after apply)
        + instance_state               = (known after apply)
        + instance_type                = "t2.small"
        + ipv6_address_count           = (known after apply)
        + ipv6_addresses               = (known after apply)
        + key_name                     = (known after apply)
        + outpost_arn                  = (known after apply)
        + password_data                = (known after apply)
        + placement_group              = (known after apply)
        + primary_network_interface_id = (known after apply)
        + private_dns                  = (known after apply)
        + private_ip                   = (known after apply)
        + public_dns                   = (known after apply)
        + public_ip                    = (known after apply)
        + secondary_private_ips        = (known after apply)
        + security_groups              = (known after apply)
        + source_dest_check            = true
        + subnet_id                    = (known after apply)
        + tenancy                      = (known after apply)
        + volume_tags                  = (known after apply)
        + vpc_security_group_ids       = (known after apply)
  
        + ebs_block_device {
            + delete_on_termination = (known after apply)
            + device_name           = (known after apply)
            + encrypted             = (known after apply)
            + iops                  = (known after apply)
            + kms_key_id            = (known after apply)
            + snapshot_id           = (known after apply)
            + volume_id             = (known after apply)
            + volume_size           = (known after apply)
            + volume_type           = (known after apply)
          }
  
        + ephemeral_block_device {
            + device_name  = (known after apply)
            + no_device    = (known after apply)
            + virtual_name = (known after apply)
          }
  
        + metadata_options {
            + http_endpoint               = (known after apply)
            + http_put_response_hop_limit = (known after apply)
            + http_tokens                 = (known after apply)
          }
  
        + network_interface {
            + delete_on_termination = (known after apply)
            + device_index          = (known after apply)
            + network_interface_id  = (known after apply)
          }
  
        + root_block_device {
            + delete_on_termination = (known after apply)
            + device_name           = (known after apply)
            + encrypted             = (known after apply)
            + iops                  = (known after apply)
            + kms_key_id            = (known after apply)
            + volume_id             = (known after apply)
            + volume_size           = (known after apply)
            + volume_type           = (known after apply)
          }
      }
  
    # aws_instance.web[0] will be created
    + resource "aws_instance" "web" {
        + ami                          = "ami-0b289b3e97908e84e"
        + arn                          = (known after apply)
        + associate_public_ip_address  = (known after apply)
        + availability_zone            = (known after apply)
        + cpu_core_count               = (known after apply)
        + cpu_threads_per_core         = (known after apply)
        + get_password_data            = false
        + host_id                      = (known after apply)
        + id                           = (known after apply)
        + instance_state               = (known after apply)
        + instance_type                = "t2.small"
        + ipv6_address_count           = (known after apply)
        + ipv6_addresses               = (known after apply)
        + key_name                     = (known after apply)
        + outpost_arn                  = (known after apply)
        + password_data                = (known after apply)
        + placement_group              = (known after apply)
        + primary_network_interface_id = (known after apply)
        + private_dns                  = (known after apply)
        + private_ip                   = (known after apply)
        + public_dns                   = (known after apply)
        + public_ip                    = (known after apply)
        + secondary_private_ips        = (known after apply)
        + security_groups              = (known after apply)
        + source_dest_check            = true
        + subnet_id                    = (known after apply)
        + tags                         = {
            + "Name" = "My_first_instance"
          }
        + tenancy                      = (known after apply)
        + volume_tags                  = (known after apply)
        + vpc_security_group_ids       = (known after apply)
  
        + ebs_block_device {
            + delete_on_termination = (known after apply)
            + device_name           = (known after apply)
            + encrypted             = (known after apply)
            + iops                  = (known after apply)
            + kms_key_id            = (known after apply)
            + snapshot_id           = (known after apply)
            + volume_id             = (known after apply)
            + volume_size           = (known after apply)
            + volume_type           = (known after apply)
          }
  
        + ephemeral_block_device {
            + device_name  = (known after apply)
            + no_device    = (known after apply)
            + virtual_name = (known after apply)
          }
  
        + metadata_options {
            + http_endpoint               = (known after apply)
            + http_put_response_hop_limit = (known after apply)
            + http_tokens                 = (known after apply)
          }
  
        + network_interface {
            + delete_on_termination = (known after apply)
            + device_index          = (known after apply)
            + network_interface_id  = (known after apply)
          }
  
        + root_block_device {
            + delete_on_termination = (known after apply)
            + device_name           = (known after apply)
            + encrypted             = (known after apply)
            + iops                  = (known after apply)
            + kms_key_id            = (known after apply)
            + volume_id             = (known after apply)
            + volume_size           = (known after apply)
            + volume_type           = (known after apply)
          }
      }
  
    # aws_instance.web[1] will be created
    + resource "aws_instance" "web" {
        + ami                          = "ami-0b289b3e97908e84e"
        + arn                          = (known after apply)
        + associate_public_ip_address  = (known after apply)
        + availability_zone            = (known after apply)
        + cpu_core_count               = (known after apply)
        + cpu_threads_per_core         = (known after apply)
        + get_password_data            = false
        + host_id                      = (known after apply)
        + id                           = (known after apply)
        + instance_state               = (known after apply)
        + instance_type                = "t2.small"
        + ipv6_address_count           = (known after apply)
        + ipv6_addresses               = (known after apply)
        + key_name                     = (known after apply)
        + outpost_arn                  = (known after apply)
        + password_data                = (known after apply)
        + placement_group              = (known after apply)
        + primary_network_interface_id = (known after apply)
        + private_dns                  = (known after apply)
        + private_ip                   = (known after apply)
        + public_dns                   = (known after apply)
        + public_ip                    = (known after apply)
        + secondary_private_ips        = (known after apply)
        + security_groups              = (known after apply)
        + source_dest_check            = true
        + subnet_id                    = (known after apply)
        + tags                         = {
            + "Name" = "My_first_instance"
          }
        + tenancy                      = (known after apply)
        + volume_tags                  = (known after apply)
        + vpc_security_group_ids       = (known after apply)
  
        + ebs_block_device {
            + delete_on_termination = (known after apply)
            + device_name           = (known after apply)
            + encrypted             = (known after apply)
            + iops                  = (known after apply)
            + kms_key_id            = (known after apply)
            + snapshot_id           = (known after apply)
            + volume_id             = (known after apply)
            + volume_size           = (known after apply)
            + volume_type           = (known after apply)
          }
  
        + ephemeral_block_device {
            + device_name  = (known after apply)
            + no_device    = (known after apply)
            + virtual_name = (known after apply)
          }
  
        + metadata_options {
            + http_endpoint               = (known after apply)
            + http_put_response_hop_limit = (known after apply)
            + http_tokens                 = (known after apply)
          }
  
        + network_interface {
            + delete_on_termination = (known after apply)
            + device_index          = (known after apply)
            + network_interface_id  = (known after apply)
          }
  
        + root_block_device {
            + delete_on_termination = (known after apply)
            + device_name           = (known after apply)
            + encrypted             = (known after apply)
            + iops                  = (known after apply)
            + kms_key_id            = (known after apply)
            + volume_id             = (known after apply)
            + volume_size           = (known after apply)
            + volume_type           = (known after apply)
          }
      }
  
  Plan: 4 to add, 0 to change, 0 to destroy.
  
  Changes to Outputs:
    + account_id  = "536867559799"
    + caller_user = "AIDAXZ76M3F35JWM6IYMP"
    + region_name = "us-east-2"
  
  ------------------------------------------------------------------------
  
  Note: You didn't specify an "-out" parameter to save this plan, so Terraform
  can't guarantee that exactly these actions will be performed if
  "terraform apply" is subsequently run.
  
  Releasing state lock. This may take a few moments...
  ```

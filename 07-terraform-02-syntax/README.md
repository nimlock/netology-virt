# Домашнее задание к занятию "7.2. Облачные провайдеры и синтаксис Терраформ."

## Модуль 7. Облачная инфраструктура. Terraform

### Студент: Иван Жиляев

>Зачастую разбираться в новых инструментах гораздо интересней понимая то, как они работают изнутри. 
>Поэтому в рамках первого *необязательного* задания предлагается завести свою учетную запись в AWS (Amazon Web Services).

## Задача 1. Регистрация в aws и знакомство с основами (необязательно, но крайне желательно).

>Остальные задания можно будет выполнять и без этого аккаунта, но с ним можно будет увидеть полный цикл процессов. 
>
>AWS предоставляет достаточно много бесплатных ресурсов в первых год после регистрации, подробно описано [здесь](https://aws.amazon.com/free/).
>1. Создайте аккаут aws.
>1. Установите c aws-cli https://aws.amazon.com/cli/.
>1. Выполните первичную настройку aws-sli https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html.
>1. Создайте IAM политику для терраформа c правами
>    * AmazonEC2FullAccess
>    * AmazonS3FullAccess
>    * AmazonDynamoDBFullAccess
>    * AmazonRDSFullAccess
>    * CloudWatchFullAccess
>    * IAMFullAccess
>1. Добавьте переменные окружения 
>    ```
>    export AWS_ACCESS_KEY_ID=(your access key id)
>    export AWS_SECRET_ACCESS_KEY=(your secret access key)
>    ```
>1. Создайте, остановите и удалите ec2 инстанс (любой с пометкой `free tier`) через веб интерфейс. 
>
>В виде результата задания приложите вывод команды `aws configure list`.

С созданием аккаунта AWS проблем не возникло.

Установка `aws-cli` сперва была выполнена через pip, однако в нём оказалась устаревшая мажорная версия 1. Так что установка была выполнена по рекомендации с оф.сайта:

```
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

В рамках первичной настройки AWS были произведены следующие действия:

- через web-консоль и сервис IAM была создана группа `automatization_group`, для которой были заданы указанные в задании полномочия

- через web-консоль и сервис IAM была создана сервисная учётная запись `terraform` с доступом через `Programmatic access`; запись была добавлена в группу `automatization_group`

- командой `aws configure` в профиль default были записаны реквизиты доступа созданного пользователя, а также желаемый регион и формат вывода

Вот вывод команды `aws configure list`:

```
ivan@kubang:~/study/netology-virt/07-terraform-02-syntax$ aws configure list
      Name                    Value             Type    Location
      ----                    -----             ----    --------
   profile                <not set>             None    None
access_key     ****************TRVH shared-credentials-file    
secret_key     ****************h9Z9 shared-credentials-file    
    region                us-east-2      config-file    ~/.aws/config
```

С запуском и остановкой инстанса через web-консоль проблем не возникло.

## Задача 2. Созданием ec2 через терраформ. 

>1. В каталоге `terraform` вашего основного репозитория, который был создан в начале курсе, создайте файл `main.tf` и `versions.tf`.
>1. Зарегистрируйте провайдер для [aws](https://registry.terraform.io/providers/hashicorp/aws/latest/docs). В файл `main.tf` добавьте
>блок `provider`, а в `versions.tf` блок `terraform` с вложенным блоком `required_providers`. Укажите любой выбранный вами регион 
>внутри блока `provider`.
>1. Внимание! В гит репозиторий нельзя пушить ваши личные ключи доступа к аккаунта. Поэтому в предыдущем задании мы указывали
>их в виде переменных окружения. 
>1. В файле `main.tf` воспользуйтесь блоком `data "aws_ami` для поиска ami образа последнего Ubuntu.  
>1. В файле `main.tf` создайте рессурс [ec2 instance](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance).
>Постарайтесь указать как можно больше параметров для его определения. Минимальный набор параметров указан в первом блоке 
>`Example Usage`, но желательно, указать большее количество параметров. 
>1. Добавьте data-блоки `aws_caller_identity` и `aws_region`.
>1. В файл `outputs.tf` поместить блоки `output` с данными об используемых в данный момент: 
>    * AWS account ID,
>    * AWS user ID,
>    * AWS регион, который используется в данный момент, 
>    * Приватный IP ec2 инстансы,
>    * Идентификатор подсети в которой создан инстанс.  
>1. Если вы выполнили первый пункт, то добейтесь того, что бы команда `terraform plan` выполнялась без ошибок. 
>
>
>В качестве результата задания предоставьте:
>1. Ответ на вопрос: при помощи какого инструмента (из разобранных на прошлом занятии) можно создать свой образ ami?
>1. Ссылку на репозиторий с исходной конфигурацией терраформа.  
 
Свой образ ami можно создать с помощью [HashiCorp Packer](https://www.packer.io/).

Репозиторий terraform-а я расположил в этом же git-репозитории на уровень выше в папке [terraform_repository](../terraform_repository).

Вывод команды `terraform plan`:

```
ivan@kubang:~/study/netology-virt/terraform$ terraform plan    

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.aws_instance.current will be read during apply
  # (config refers to values not yet known)
 <= data "aws_instance" "current"  {
      + ami                         = (known after apply)
      + arn                         = (known after apply)
      + associate_public_ip_address = (known after apply)
      + availability_zone           = (known after apply)
      + credit_specification        = (known after apply)
      + disable_api_termination     = (known after apply)
      + ebs_block_device            = (known after apply)
      + ebs_optimized               = (known after apply)
      + ephemeral_block_device      = (known after apply)
      + host_id                     = (known after apply)
      + iam_instance_profile        = (known after apply)
      + id                          = (known after apply)
      + instance_id                 = (known after apply)
      + instance_state              = (known after apply)
      + instance_tags               = (known after apply)
      + instance_type               = (known after apply)
      + key_name                    = (known after apply)
      + metadata_options            = (known after apply)
      + monitoring                  = (known after apply)
      + network_interface_id        = (known after apply)
      + outpost_arn                 = (known after apply)
      + password_data               = (known after apply)
      + placement_group             = (known after apply)
      + private_dns                 = (known after apply)
      + private_ip                  = (known after apply)
      + public_dns                  = (known after apply)
      + public_ip                   = (known after apply)
      + root_block_device           = (known after apply)
      + secondary_private_ips       = (known after apply)
      + security_groups             = (known after apply)
      + source_dest_check           = (known after apply)
      + subnet_id                   = (known after apply)
      + tags                        = (known after apply)
      + tenancy                     = (known after apply)
      + user_data                   = (known after apply)
      + user_data_base64            = (known after apply)
      + vpc_security_group_ids      = (known after apply)
    }

  # aws_instance.web will be created
  + resource "aws_instance" "web" {
      + ami                          = "ami-01112dc69308aab5e"
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

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + account_id          = "536867559799"
  + caller_user         = "AIDAXZ76M3F35JWM6IYMP"
  + instance_private_ip = (known after apply)
  + instance_subnet_id  = (known after apply)
  + region_name         = "us-east-2"

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

terraform {
  backend "s3" {
    bucket = "ivan-07-terraform-bucket"
    key    = "terraform-remote-config/terraform.tfstate"
    region = "us-east-2"
    dynamodb_table = "terraform-state-locking"
  }
}

# data "terraform_remote_state" "cloudinfra" {
#   backend = "s3"
#   config = {
#     bucket = "ivan-07-terraform-bucket"
#     key    = "terraform-remote-config/terraform.tfstate"
#     region = "us-east-2"
#   }
# }

provider "aws" {
    region = "us-east-2"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

locals {
  dict_of_instance_types = {
    stage = "t2.micro"
    prod = "t2.small"
  }
}

locals {
  dict_of_instance_count = {
    stage = 1
    prod = 2
  }
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = local.dict_of_instance_types[terraform.workspace]
  count = local.dict_of_instance_count[terraform.workspace]

  tags = {
    Name = "My_first_instance"
  }

  lifecycle {
    create_before_destroy = true
  }

}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

# data "aws_instance" "current" {
#   instance_id = aws_instance.web.id
# }

locals {
  instances = {
    "t2.micro" = data.aws_ami.ubuntu.id
    "t2.small" = data.aws_ami.ubuntu.id
  }
}

resource "aws_instance" "api" {
  for_each = local.instances

  ami           = each.value
  instance_type = each.key
}

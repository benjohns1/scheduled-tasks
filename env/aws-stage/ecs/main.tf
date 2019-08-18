terraform {
  required_version = ">= 0.12"
}

variable "tags" {
  type = map(string)
}

variable "prefix" {
  description = "Prefix for resource names"
  type = string
}

variable "logregion" {
  type = string
}

variable "host_application_port" {
  type = number
}

variable "host_webapp_port" {
  type = number
}

variable "aws_ec2_public_key_name" {
  type = string
}

variable "aws_ec2_public_key" {
  type = string
}

variable "container_env" {
  type = map(string)
  default = {
    "APPLICATION_PORT" = "3000",
    "AUTH0_DOMAIN" = "",
    "AUTH0_API_IDENTIFIER" = "",
    "AUTH0_API_SECRET" = "",
    "POSTGRES_HOST" = "",
    "POSTGRES_PORT" = "5432",
    "POSTGRES_DB" = "",
    "POSTGRES_USER" = "",
    "POSTGRES_PASSWORD" = "",
    "DBCONN_MAXRETRYATTEMPTS" = "20",
    "DBCONN_RETRYSLEEPSECONDS" = "3",
    "WEBAPP_PORT" = "80",
    "AUTH0_WEBAPP_CLIENT_ID" = "",
    "AUTH0_ANON_CLIENT_ID" = "",
    "AUTH0_ANON_CLIENT_SECRET" = "",
    "ENABLE_E2E_DEV_LOGIN" = "false",
    "AUTH0_E2E_DEV_CLIENT_ID" = "",
    "AUTH0_E2E_DEV_CLIENT_SUBJECT" = "",
    "AUTH0_E2E_DEV_CLIENT_SECRET" = ""
  }
}

locals {
  logname_services = "/ecs/${var.prefix}-services"
  logname_webapp = "/ecs/${var.prefix}-webapp"
}

resource "aws_ecs_cluster" "ecs_cluster" {
  name = "${var.prefix}-ecs-cluster"
  tags = var.tags
}

resource "aws_key_pair" "ecs_keypair" {
  key_name = var.aws_ec2_public_key_name
  public_key = var.aws_ec2_public_key
}

resource "aws_launch_configuration" "ecs_launch_config" {
  name = "${var.prefix}-launch-config"
  image_id = "ami-0e5e051fd0b505db6"
  instance_type = "t3a.micro"
  security_groups = [aws_security_group.nsg_task.id]
  iam_instance_profile = aws_iam_instance_profile.ecs_instance_profile.id
  key_name = var.aws_ec2_public_key_name
  lifecycle {
    create_before_destroy = true
  }
  root_block_device {
    volume_type = "gp2"
    volume_size = 30
  }
  user_data = <<USER_DATA
#!/bin/bash
echo ECS_CLUSTER=${aws_ecs_cluster.ecs_cluster.name} >> /etc/ecs/ecs.config
USER_DATA
}

resource "aws_autoscaling_group" "ecs_autoscaling_group" {
  name = "${var.prefix}-autoscaling-group"
  max_size = 2
  min_size = 1
  desired_capacity = 1
  vpc_zone_identifier = data.aws_subnet_ids.all.ids
  launch_configuration = aws_launch_configuration.ecs_launch_config.name
  tag {
    key = "Name"
    value = "${var.prefix}-ec2-instance"
    propagate_at_launch = true
  }
}

data "aws_instance" "ec2_instance" {
  depends_on = ["aws_autoscaling_group.ecs_autoscaling_group"]
  filter {
    name = "tag:Name"
    values = ["${var.prefix}-ec2-instance"]
  }
}

resource "aws_eip" "ec2_eip" {
  vpc = true
  instance = data.aws_instance.ec2_instance.id
  associate_with_private_ip = data.aws_instance.ec2_instance.private_ip
}

resource "aws_ecs_task_definition" "tasks" {
  family = "${var.prefix}-tasks"
  container_definitions = <<CONTAINER_DEFS
[
  ${templatefile("${path.module}/services_container_definition.json", {
    name = "${var.prefix}-services",
    logname = local.logname_services,
    logregion = var.logregion,
    host_application_port = var.host_application_port,
    env = var.container_env
  })},
  ${templatefile("${path.module}/webapp_container_definition.json", {
    name = "${var.prefix}-webapp",
    logname = local.logname_webapp,
    logregion = var.logregion,
    application_host = "localhost",
    host_webapp_port = var.host_webapp_port,
    env = var.container_env
  })}
]
CONTAINER_DEFS
  requires_compatibilities = ["EC2"]
  cpu = "256"
  memory = "512"
  network_mode = "host"
  execution_role_arn = aws_iam_role.ecs_task_role.arn
  tags = var.tags
  /*lifecycle {
    ignore_changes = [
      container_definitions
    ]
  }*/
}

resource "aws_ecs_service" "ecs_service" {
  name = "${var.prefix}-tasks"
  cluster = aws_ecs_cluster.ecs_cluster.id
  task_definition = aws_ecs_task_definition.tasks.arn
  desired_count = 1
  tags = var.tags
  enable_ecs_managed_tags = true
}

resource "aws_cloudwatch_log_group" "services_logs" {
  name = local.logname_services
  retention_in_days = "7"
  tags = var.tags
}

resource "aws_cloudwatch_log_group" "webapp_logs" {
  name = local.logname_webapp
  retention_in_days = "7"
  tags = var.tags
}

output "host_webapp_port" {
  value = var.container_env["WEBAPP_PORT"]
}

output "host_public_ip_addr" {
  value = aws_eip.ec2_eip.public_ip
}

output "host_private_ip_addr" {
  value = data.aws_instance.ec2_instance.private_ip
}

output "host_public_dns" {
  value = aws_route53_record.subdomain.name
}

terraform {
  required_version = ">= 0.12"
}

variable "cluster_id" {
    description = "The ECS cluster ID"
    type = string
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
    "DBCONN_RETRYSLEEPSECONDS" = "3"
  }
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_execution_role" {
    name = "${var.prefix}-ecs-execution-role"
    assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
}

resource "aws_iam_policy_attachment" "ecs_execution_policy" {
    name = "${var.prefix}-ecs-execution-policy"
    roles = [aws_iam_role.ecs_execution_role.name]
    policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

locals {
  logname = "/ecs/${var.prefix}-services"
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnet_ids" "all" {
  vpc_id = data.aws_vpc.default.id
}

resource "aws_security_group" "nsg_task" {
  name = "${var.prefix}-task"
  description = "${var.prefix} services ports"
  vpc_id = data.aws_vpc.default.id
  tags = var.tags
}

resource "aws_security_group_rule" "nsg_task_ingress_rule" {
  security_group_id        = "${aws_security_group.nsg_task.id}"
  type                     = "ingress"
  protocol                 = "tcp"
  cidr_blocks              = ["0.0.0.0/0"]
  from_port                = "${var.host_application_port}"
  to_port                  = "${var.host_application_port}"
}

resource "aws_security_group_rule" "nsg_task_egress_rule" {
  security_group_id        = "${aws_security_group.nsg_task.id}"
  type                     = "egress"
  protocol                 = "-1"
  cidr_blocks              = ["0.0.0.0/0"]
  from_port                = "0"
  to_port                  = "0"
}

resource "aws_ecs_task_definition" "services" {
    family = "${var.prefix}-services"
    container_definitions = templatefile("${path.module}/services_container_definition.json", {
      prefix = var.prefix,
      logname = local.logname,
      logregion = var.logregion,
      host_application_port = var.host_application_port,
      env = var.container_env
    })
    requires_compatibilities = ["FARGATE"]
    cpu = "256"
    memory = "512"
    network_mode = "awsvpc"
    execution_role_arn = aws_iam_role.ecs_execution_role.arn
    tags = var.tags
}

resource "aws_ecs_service" "services" {
    name = "${var.prefix}-services"
    cluster = var.cluster_id
    task_definition = aws_ecs_task_definition.services.arn
    desired_count = 1
    launch_type = "FARGATE"
    deployment_maximum_percent = 100
    deployment_minimum_healthy_percent = 50
    network_configuration {
        security_groups = [aws_security_group.nsg_task.id]
        subnets = data.aws_subnet_ids.all.ids
        assign_public_ip = true
    }
    tags = var.tags
    enable_ecs_managed_tags = true

    lifecycle {
      ignore_changes = [task_definition]
    }
}

resource "aws_cloudwatch_log_group" "logs" {
  name = local.logname
  retention_in_days = "7"
  tags = var.tags
}
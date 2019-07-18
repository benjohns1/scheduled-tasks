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
  logname = "/fargate/service/${var.prefix}-services"
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnet_ids" "all" {
  vpc_id = data.aws_vpc.default.id
}

data "aws_security_group" "default" {
  vpc_id = data.aws_vpc.default.id
  name = "default"
}

resource "aws_ecs_task_definition" "services" {
    family = "services"
    container_definitions = templatefile("${path.module}/task_definition.json", {
      logname = local.logname,
      logregion = var.logregion
    })
    requires_compatibilities = ["FARGATE"]
    cpu = "256"
    memory = "512"
    network_mode = "awsvpc"
    execution_role_arn = aws_iam_role.ecs_execution_role.arn
    task_role_arn = aws_iam_role.ecs_execution_role.arn
    tags = var.tags
}

resource "aws_ecs_service" "services" {
    name = "services"
    cluster = var.cluster_id
    task_definition = aws_ecs_task_definition.services.arn
    desired_count = 1
    launch_type = "FARGATE"
    deployment_maximum_percent = 100
    deployment_minimum_healthy_percent = 0
    network_configuration {
        security_groups = [data.aws_security_group.default.id]
        subnets = data.aws_subnet_ids.all.ids
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
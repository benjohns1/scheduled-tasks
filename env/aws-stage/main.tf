terraform {
  required_version = ">= 0.12"
}

locals {
  prefix = "${var.application_name}-${var.env}"
  db_name = "${local.prefix}-db-${var.postgres_db_name}"
  tags = {
    Environment = var.env
  }
}

provider "aws" {
  version = "~> 2.0"
  region = var.aws_region
  profile = var.aws_profile
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

module "db" {

  source = "terraform-aws-modules/rds/aws"
  version = "~> 2.0"

  identifier = local.db_name

  engine            = "postgres"
  engine_version    = "10.6"
  instance_class    = "db.t2.micro"
  allocated_storage = 20
  storage_encrypted = false
  publicly_accessible = true

  name = var.postgres_db_name
  username = var.postgres_db_user
  password = var.postgres_db_password
  port     = var.postgres_db_port

  vpc_security_group_ids = [data.aws_security_group.default.id]
  performance_insights_enabled = false
  performance_insights_retention_period = 0

  maintenance_window = "Mon:00:00-Mon:03:00"
  backup_window      = "03:00-06:00"

  # disable backups to create DB faster
  backup_retention_period = 0

  tags = local.tags

  # DB subnet group
  subnet_ids = data.aws_subnet_ids.all.ids

  # DB parameter group
  family = "postgres10"

  # DB option group
  major_engine_version = "10.6"

  # Snapshot name upon DB deletion
  final_snapshot_identifier = local.db_name

  # Database Deletion Protection
  deletion_protection = false
}

resource "aws_ecs_cluster" "ecs_cluster" {
  name = "${local.prefix}-ecs-cluster"
  tags = local.tags
}

module "services" {
  source = "./services"
  cluster_id = aws_ecs_cluster.ecs_cluster.id
  tags = local.tags
  prefix = local.prefix
  logregion = var.aws_region
}
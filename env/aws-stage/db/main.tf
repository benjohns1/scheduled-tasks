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

variable "db_name" {
  type = string
}

variable "db_user" {
  type = string
}

variable "db_password" {
  type = string
}

variable "db_port" {
  default = 5432
  type = number
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnet_ids" "all" {
  vpc_id = data.aws_vpc.default.id
}

resource "aws_security_group" "nsg_db" {
  name = "${var.prefix}-db"
  description = "${var.prefix} DB ports"
  vpc_id = data.aws_vpc.default.id
  tags = var.tags
}

resource "aws_security_group_rule" "nsg_db_ingress_rule" {
  security_group_id        = "${aws_security_group.nsg_db.id}"
  type                     = "ingress"
  protocol                 = "tcp"
  cidr_blocks              = ["0.0.0.0/0"]
  from_port                = "${var.db_port}"
  to_port                  = "${var.db_port}"
}

resource "aws_security_group_rule" "nsg_db_egress_rule" {
  security_group_id        = "${aws_security_group.nsg_db.id}"
  type                     = "egress"
  protocol                 = "tcp"
  cidr_blocks              = ["0.0.0.0/0"]
  from_port                = "0"
  to_port                  = "0"
}

locals {
  identifier = "${var.prefix}-db-${var.db_name}"
}

module "aws_rds" {

  source = "terraform-aws-modules/rds/aws"
  version = "~> 2.0"

  identifier = local.identifier

  engine            = "postgres"
  engine_version    = "10.6"
  instance_class    = "db.t2.micro"
  allocated_storage = 20
  storage_encrypted = false
  publicly_accessible = true

  name = var.db_name
  username = var.db_user
  password = var.db_password
  port     = var.db_port

  vpc_security_group_ids = [aws_security_group.nsg_db.id]
  performance_insights_enabled = false
  performance_insights_retention_period = 0

  maintenance_window = "Mon:00:00-Mon:03:00"
  backup_window      = "03:00-06:00"

  # disable backups to create DB faster
  backup_retention_period = 0

  tags = var.tags

  # DB subnet group
  subnet_ids = data.aws_subnet_ids.all.ids

  # DB parameter group
  family = "postgres10"

  # DB option group
  major_engine_version = "10.6"

  # Snapshot name upon DB deletion
  final_snapshot_identifier = local.identifier

  # Database Deletion Protection
  deletion_protection = false
}

output "db_instance_address" {
  description = "DB connection address"
  value       = module.aws_rds.this_db_instance_address
}
terraform {
  required_version = ">= 0.12"
}

locals {
  prefix = "${var.application_name}-${var.env}"
  tags = {
    Environment = var.env
  }
}

provider "aws" {
  version = "~> 2.0"
  region = var.aws_region
  profile = var.aws_profile
}

module "db" {
  source = "./db"
  prefix = local.prefix
  db_name = var.postgres_db_name
  tags = local.tags
  db_port = var.postgres_db_port
  db_user = var.postgres_db_user
  db_password = var.postgres_db_password
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
  host_application_port = var.application_port
  container_env = {
    "APPLICATION_PORT" = var.application_port,
    "AUTH0_DOMAIN" = var.auth0_domain,
    "AUTH0_API_IDENTIFIER" = var.auth0_api_identifier,
    "AUTH0_API_SECRET" = var.auth0_api_secret,
    "POSTGRES_HOST" = module.db.db_instance_endpoint,
    "POSTGRES_PORT" = "${var.postgres_db_port}",
    "POSTGRES_DB" = var.postgres_db_name,
    "POSTGRES_USER" = var.postgres_db_user,
    "POSTGRES_PASSWORD" = var.postgres_db_password,
    "DBCONN_MAXRETRYATTEMPTS" = "${var.dbconn_maxretryattempts}",
    "DBCONN_RETRYSLEEPSECONDS" = "${var.dbconn_retrysleepseconds}"
  }
}
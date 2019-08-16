terraform {
  required_version = ">= 0.12"
}

locals {
  prefix = "${var.application_name}-${var.env}"
  tags = {
    Name = var.application_name
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

/* @TODO: create vpc with public/private subnets, use load balancer for public traffic
 * https://airship.tf/getting_started/preparation.html#vpc
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "2.7.0"

  name = "${local.prefix}-vpc"
  cidr = "10.50.0.0/16"

  azs = ["us-west-2a", "us-west-2b"]
  public_subnets = ["10.50.11.0/24", "10.50.12.0/24"]
  private_subnets = ["10.50.21.0/24", "10.50.22.0/24"]

  single_nat_gateway = true

  enable_nat_gateway   = true
  enable_vpn_gateway   = false
  enable_dns_hostnames = true

  tags = local.tags

}*/

module "db" {
  source = "./db"
  prefix = local.prefix
  db_name = var.postgres_db_name
  tags = local.tags
  db_port = var.postgres_db_port
  db_user = var.postgres_db_user
  db_password = var.postgres_db_password
}

/* @TODO: create service discovery resources for container networking
resource "aws_service_discovery_private_dns_namespace" "dns_namespace" {
  name = "${local.prefix}-dns-namespace"
  vpc = data.aws_vpc.default.id
}*/

module "ecs" {
  source = "./ecs"
  tags = local.tags
  prefix = local.prefix
  logregion = var.aws_region
  host_application_port = var.application_port
  host_webapp_port = var.webapp_port
  aws_ec2_public_key_name = var.aws_ec2_public_key_name
  aws_ec2_public_key = var.aws_ec2_public_key
  container_env = {
    "APPLICATION_PORT" = var.application_port,
    "AUTH0_DOMAIN" = var.auth0_domain,
    "AUTH0_API_IDENTIFIER" = var.auth0_api_identifier,
    "AUTH0_API_SECRET" = var.auth0_api_secret,
    "POSTGRES_HOST" = module.db.db_instance_address,
    "POSTGRES_PORT" = "${var.postgres_db_port}",
    "POSTGRES_DB" = var.postgres_db_name,
    "POSTGRES_USER" = var.postgres_db_user,
    "POSTGRES_PASSWORD" = var.postgres_db_password,
    "DBCONN_MAXRETRYATTEMPTS" = "${var.dbconn_maxretryattempts}",
    "DBCONN_RETRYSLEEPSECONDS" = "${var.dbconn_retrysleepseconds}",
    "WEBAPP_PORT" = var.webapp_port,
    "AUTH0_WEBAPP_CLIENT_ID" = var.auth0_webapp_client_id,
    "AUTH0_ANON_CLIENT_ID" = var.auth0_anon_client_id,
    "AUTH0_ANON_CLIENT_SECRET" = var.auth0_anon_client_secret,
    "ENABLE_E2E_DEV_LOGIN" = "${var.enable_e2e_dev_login}",
    "AUTH0_E2E_DEV_CLIENT_ID" = var.auth0_e2e_dev_client_id,
    "AUTH0_E2E_DEV_CLIENT_SUBJECT" = var.auth0_e2e_dev_client_subject,
    "AUTH0_E2E_DEV_CLIENT_SECRET" = var.auth0_e2e_dev_client_secret
  }
}

output "host_webapp_port" {
  value = module.ecs.host_webapp_port
}

output "host_public_ip_addr" {
  value = module.ecs.host_public_ip_addr
}

output "host_public_dns" {
  value = module.ecs.host_public_dns
}

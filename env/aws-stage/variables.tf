variable "postgres_db_name" {
}

variable "postgres_db_user" {
}

variable "POSTGRES_PASSWORD" {
}

variable "POSTGRES_PORT" {
}

variable "env" {
    default = "prod"
}

variable "aws_region" {
    default = "us-west-2"
}

variable "aws_profile" {
    default = ""
}

variable "application_name" {
    default = "st"
}
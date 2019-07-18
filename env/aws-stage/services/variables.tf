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
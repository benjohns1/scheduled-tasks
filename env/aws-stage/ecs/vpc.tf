data "aws_vpc" "default" {
  default = true
}

data "aws_subnet_ids" "all" {
  vpc_id = data.aws_vpc.default.id
}

resource "aws_security_group" "nsg_task" {
  name = "${var.prefix}-task"
  description = "${var.prefix} services and webapp ports"
  vpc_id = data.aws_vpc.default.id
  tags = var.tags
  ingress {
    protocol                 = "tcp"
    cidr_blocks              = ["0.0.0.0/0"]
    from_port                = "${var.host_webapp_port}"
    to_port                  = "${var.host_webapp_port}"
  }
  egress {
    protocol                 = "-1"
    cidr_blocks              = ["0.0.0.0/0"]
    from_port                = "0"
    to_port                  = "0"
  }
}

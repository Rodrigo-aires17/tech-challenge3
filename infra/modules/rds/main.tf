resource "random_password" "this" {
  length  = 24
  special = false # evita caracteres que exijam escaping em connection strings
}

resource "aws_db_subnet_group" "this" {
  name       = "${var.name_prefix}-${var.identifier}-subnet-group"
  subnet_ids = var.subnet_ids
  tags       = var.tags
}

resource "aws_security_group" "this" {
  name        = "${var.name_prefix}-${var.identifier}-rds-sg"
  description = "Acesso ao RDS ${var.identifier} somente a partir da VPC/EKS"
  vpc_id      = var.vpc_id

  ingress {
    description     = "PostgreSQL a partir dos security groups permitidos"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = var.allowed_security_group_ids
    cidr_blocks     = length(var.allowed_security_group_ids) == 0 ? [var.vpc_cidr_block] : null
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, { Name = "${var.name_prefix}-${var.identifier}-rds-sg" })
}

resource "aws_db_instance" "this" {
  identifier     = "${var.name_prefix}-${var.identifier}"
  engine         = "postgres"
  engine_version = var.engine_version
  instance_class = var.instance_class

  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.allocated_storage * 5
  storage_encrypted     = true

  db_name  = "${var.identifier}_db"
  username = "${var.identifier}_admin"
  password = random_password.this.result
  port     = 5432

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.this.id]

  multi_az                = false
  publicly_accessible     = false
  backup_retention_period = 7
  skip_final_snapshot     = true
  deletion_protection     = false
  apply_immediately       = true

  tags = var.tags
}

resource "aws_secretsmanager_secret" "this" {
  name = "${var.name_prefix}-${var.identifier}-db"
  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "this" {
  secret_id = aws_secretsmanager_secret.this.id
  secret_string = jsonencode({
    username = aws_db_instance.this.username
    password = random_password.this.result
    host     = aws_db_instance.this.address
    port     = aws_db_instance.this.port
    dbname   = aws_db_instance.this.db_name
  })
}

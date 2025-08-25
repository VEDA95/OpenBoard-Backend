data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/loader/main.go",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  url = "postgres://Testuser:GNX-0093@localhost:5432/open_board_api?sslmode=disable"
  dev = "postgres://Testuser:GNX-0093@localhost:5432/open_board_api?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

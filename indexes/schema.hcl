schema "main" {}

table "stars" {
  schema = schema.main

  column "id" {
    type = text
    null = false
  }

  column "designation" {
    type = text
    null = false
  }

  column "x" {
    type = float
    null = false
  }

  column "y" {
    type = float
    null = false
  }

  column "ra" {
    type = float
    null = false
  }

  column "dec" {
    type = float
    null = false
  }

  column "intensity" {
    type = float
    null = false
  }

  column "pixel" {
    type = integer
    null = false
  }

  primary_key {
    columns = [column.designation]
  }

  index "idx_pixel" {
    columns = [column.pixel]
  }
}
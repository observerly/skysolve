# **************************************************************************************

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

# **************************************************************************************

// Define an environment named "512", which corresponds to HEALPix NSides=512.
env "512" {
  // Declare where the schema definition resides.
  // Also supported: ["file://multi.hcl", "file://schema.hcl"].
  src = "file://indexes/schema.hcl"

  // Define the URL of the database which is managed
  // in this environment.
  url = "sqlite://indexes/512/stars.db.sqlite"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "sqlite://indexes/512/stars.db.sqlite"

  migration {
    // URL where the migration directory resides.
    dir = "file://migrations"
  }
}

# **************************************************************************************

// Define an environment named "1024", which corresponds to HEALPix NSides=1024.
env "1024" {
  // Declare where the schema definition resides.
  // Also supported: ["file://multi.hcl", "file://schema.hcl"].
  src = "file://indexes/schema.hcl"

  // Define the URL of the database which is managed
  // in this environment.
  url = "sqlite://indexes/1024/stars.db.sqlite"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "sqlite://indexes/1024/stars.db.sqlite"

  migration {
    // URL where the migration directory resides.
    dir = "file://migrations"
  }
}

# **************************************************************************************
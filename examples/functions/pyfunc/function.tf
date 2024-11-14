terraform {
  required_providers {
    pyfunc = {
      source = "abex.dev/abex/pyfunc"
    }
  }
}

provider "pyfunc" {}

output "pyeval" {
  value = trimspace(provider::pyfunc::pyeval("import sys; print(sys.executable)").stdout)
}

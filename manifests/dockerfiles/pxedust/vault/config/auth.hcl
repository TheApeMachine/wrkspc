auth {
  ldap {
    description = "LDAP Auth backend config"
    authconfig {
      binddn   = "CN=SamE,CN=Users,DC=test,DC=local"
      bindpass = "z"
      url      = "ldap://10.255.0.30"
      userdn   = "CN=Users,DC=test,DC=local"
    }
    group "customer" {
      options {
        policies = ["customer-policy"]
      }
    }
    mountconfig {
      default_lease_ttl = "1h"
      max_lease_ttl     = "24h"
    }
  }
  github {
    authconfig {
      organization = "<your org>"
    }
    team "backend" {
      options {
        policies = ["developer-backend-policy"]
      }
    }
  }
}
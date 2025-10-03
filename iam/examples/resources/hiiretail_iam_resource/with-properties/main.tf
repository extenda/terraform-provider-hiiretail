terraform {
  required_version = ">= 1.0"
  
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 1.0"
    }
  }
}

provider "hiiretail" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  token_url     = var.token_url
}

# Store resource with comprehensive properties
resource "hiiretail_iam_resource" "flagship_store" {
  id   = "store:flagship:nyc"
  name = "Flagship Store - New York City"
  
  props = jsonencode({
    # Location information
    location = {
      address      = "123 Broadway"
      city         = "New York"
      state        = "NY"
      zip          = "10001"
      coordinates = {
        lat = 40.7589
        lng = -73.9851
      }
    }
    
    # Store details
    store_details = {
      square_footage = 15000
      floors         = 3
      opened_date    = "2020-03-15"
      manager        = "sarah.johnson@company.com"
      phone          = "+1-212-555-0123"
    }
    
    # Departments
    departments = [
      "electronics",
      "clothing", 
      "home-garden",
      "beauty",
      "books"
    ]
    
    # Operating hours
    hours = {
      monday    = "09:00-21:00"
      tuesday   = "09:00-21:00"
      wednesday = "09:00-21:00"
      thursday  = "09:00-21:00"
      friday    = "09:00-22:00"
      saturday  = "09:00-22:00"
      sunday    = "10:00-20:00"
    }
    
    # Features and capabilities
    features = {
      wifi_available     = true
      parking_available  = true
      wheelchair_access  = true
      click_and_collect  = true
      returns_accepted   = true
      gift_wrapping     = true
    }
    
    # Staff information
    staff = {
      total_employees = 85
      managers        = 12
      seasonal_staff  = 25
    }
    
    # Performance metrics (example)
    metrics = {
      daily_target_sales    = 50000
      annual_revenue_target = 18000000
      customer_rating      = 4.7
      nps_score           = 72
    }
  })
}

# Department resource with detailed properties
resource "hiiretail_iam_resource" "electronics_department" {
  id   = "dept:flagship:electronics"
  name = "Electronics Department - Flagship Store"
  
  props = jsonencode({
    # Parent store reference
    store_id = hiiretail_iam_resource.flagship_store.id
    store_name = hiiretail_iam_resource.flagship_store.name
    
    # Department details
    department_info = {
      manager       = "mike.chen@company.com"
      assistant_mgr = "lisa.wong@company.com"
      floor         = 2
      square_footage = 3500
      employee_count = 15
    }
    
    # Product categories
    categories = {
      mobile = {
        brands = ["Apple", "Samsung", "Google", "OnePlus"]
        budget = 2000000
      }
      computers = {
        brands = ["Apple", "Dell", "HP", "Lenovo", "Microsoft"]
        budget = 1500000
      }
      gaming = {
        brands = ["Sony", "Microsoft", "Nintendo", "Razer"]
        budget = 800000
      }
      accessories = {
        types = ["cases", "chargers", "headphones", "cables"]
        budget = 300000
      }
    }
    
    # Inventory management
    inventory = {
      reorder_point    = 100
      max_stock_level  = 1000
      auto_reorder     = true
      supplier_count   = 25
    }
    
    # Sales targets
    targets = {
      monthly_revenue = 750000
      quarterly_units = 5000
      margin_target   = 0.25
    }
    
    # Customer service
    services = {
      tech_support     = true
      setup_service    = true
      warranty_claims  = true
      trade_in_program = true
    }
  })
}

# POS Terminal resource with configuration
resource "hiiretail_iam_resource" "pos_terminal_01" {
  id   = "pos:flagship:terminal-01"
  name = "POS Terminal 01 - Electronics Department"
  
  props = jsonencode({
    # Hardware details
    hardware = {
      model          = "VeriFone MX925"
      serial_number  = "VF925-NYC-001"
      install_date   = "2020-03-20"
      warranty_until = "2025-03-20"
    }
    
    # Location within store
    location = {
      store_id    = hiiretail_iam_resource.flagship_store.id
      department  = hiiretail_iam_resource.electronics_department.id
      floor       = 2
      section     = "checkout-alpha"
      position    = 1
    }
    
    # Configuration
    settings = {
      auto_update    = true
      receipt_printer = true
      card_reader    = "chip-and-pin"
      nfc_enabled    = true
      cash_drawer    = true
    }
    
    # Payment methods supported
    payment_methods = [
      "cash",
      "credit-card",
      "debit-card", 
      "nfc-payment",
      "gift-card",
      "store-credit"
    ]
    
    # Network configuration
    network = {
      ip_address = "192.168.1.101"
      mac_address = "00:1A:2B:3C:4D:5E"
      wifi_enabled = false
      ethernet_port = "A1"
    }
    
    # Security
    security = {
      encryption_enabled = true
      pin_required      = true
      timeout_minutes   = 15
      last_security_update = "2024-01-15"
    }
  })
}

# Application resource for inventory management
resource "hiiretail_iam_resource" "inventory_app" {
  id   = "app:inventory-management"
  name = "Inventory Management System"
  
  props = jsonencode({
    # Application details
    application = {
      version     = "2.1.4"
      environment = "production"
      deployed_at = "2024-01-10T14:30:00Z"
      vendor      = "RetailSoft Inc"
    }
    
    # Infrastructure
    infrastructure = {
      hosting_provider = "AWS"
      region          = "us-east-1"
      instance_type   = "t3.large"
      replicas        = 3
      load_balancer   = true
    }
    
    # Database
    database = {
      type     = "PostgreSQL"
      version  = "14.9"
      size_gb  = 500
      backups  = "daily"
      encryption = true
    }
    
    # API endpoints
    endpoints = [
      "https://api.company.com/inventory/v1",
      "https://api.company.com/inventory/health",
      "https://api.company.com/inventory/metrics"
    ]
    
    # Features
    features = {
      real_time_sync     = true
      batch_import       = true
      audit_trail        = true
      reporting          = true
      mobile_app         = true
      barcode_scanning   = true
    }
    
    # Integration
    integrations = {
      pos_systems    = ["VeriFone", "Square", "Clover"]
      accounting     = ["QuickBooks", "Xero"]
      ecommerce      = ["Shopify", "WooCommerce"]
      suppliers      = ["EDI", "API", "Manual"]
    }
    
    # Monitoring
    monitoring = {
      uptime_target    = 99.9
      response_time_ms = 200
      error_rate_max   = 0.1
      alerts_enabled   = true
    }
    
    # Security
    security = {
      ssl_enabled        = true
      api_key_auth       = true
      oauth2_enabled     = true
      rate_limiting      = true
      ip_whitelist       = ["10.0.0.0/8", "192.168.0.0/16"]
    }
  })
}
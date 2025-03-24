setup:
	@echo "=== Installing Tools ==="
	@bash ./scripts/install_tools.sh
	@echo "=== Tools Installed ==="

compose:
	@echo "=== Starting Server ==="
	@docker compose -f docker-compose.yml up --build
	@echo "=== Development Server Stopped ==="

insert:
	@echo "=== Inserting Dummy Data ==="
	@bash ./scripts/insert_data.sh
	@echo "=== Dummy Data Inserted ==="

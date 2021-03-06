PYTHON ?= python3
PIP ?= $(PYTHON) -m pip
COVERAGE ?= coverage


help: Makefile
	@echo
	@echo " Choose a command run in Silverback:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo


## config: Install dependencies.
config:
	$(PIP) install pycodestyle
	$(PIP) install coverage
	$(PIP) install flake8
	$(PIP) install flake8-comprehensions
	$(PIP) install flake8-eradicate
	$(PIP) install flake8-todo
	$(PIP) install -r requirements.txt


## lint-pycodestyle: PyCode Style Lint
lint-pycodestyle:
	@echo "\n>> ============= Pycodestyle Linting ============= <<"
	@find app -type f -name \*.py | while read file; do echo "$$file" && pycodestyle --config=./pycodestyle --first "$$file" || exit 1; done


## lint-flake8: Flake8 Lint.
lint-flake8:
	@echo "\n>> ============= Flake8 Linting ============= <<"
	@find app -type f -name \*.py | while read file; do echo "$$file" && flake8 --config=flake8.ini "$$file" || exit 1; done


## lint: Lint The Code.
lint: lint-pycodestyle lint-flake8
	@echo "\n>> ============= All linting cases passed! ============= <<"


## test: Run Test Cases.
test:
	@echo "\n>> ============= Run Test Cases ============= <<"
	$(PYTHON) manage.py test


## migration: Create DB Migration Files.
migration:
	@echo "\n>> ============= Make Migrations ============= <<"
	$(PYTHON) manage.py makemigrations


## migrate: Migrate Database.
migrate:
	@echo "\n>> ============= Migrate ============= <<"
	$(PYTHON) manage.py migrate


## run: Run Server.
run:
	@echo "\n>> ============= Run Server ============= <<"
	$(PYTHON) manage.py runserver


## coverage: Get test coverage.
coverage:
	$(COVERAGE) run --source='.' manage.py test app
	$(COVERAGE) report -m
	$(COVERAGE) html


## health: Check the application health.
health:
	$(PYTHON) manage.py health check


## create-env: Create .env file.
create-env:
	cp .env.example .env


## makemessages: Make translation files.
makemessages:
	$(PYTHON) manage.py makemessages


## gh-config: Switch to sqlite database.
gh-config: config create-env migrate
	$(PYTHON) manage.py silverback update_env DB_CONNECTION=sqlite


## outdated-pkg: Show outdated python packages
outdated-pkg:
	@echo "\n>> ============= List Outdated Packages ============= <<"
	$(PIP) list --outdated


## liteci: Run a lite CI tests.
liteci: gh-config test coverage lint outdated-pkg
	@echo "\n>> ============= All quality checks passed ============= <<"


## ci: Run all CI tests.
ci: test coverage lint outdated-pkg
	@echo "\n>> ============= All quality checks passed ============= <<"


.PHONY: help
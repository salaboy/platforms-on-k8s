# renovate: datasource=github-tags depName=jaegertracing/jaeger-operator
JAEGER_VERSION ?= v1.45.0
TOOLKIT_NAMESPACE ?= keptn-lifecycle-toolkit-system
GRAFANA_PORT_FORWARD ?= 3000

.PHONY: install
install: install-observability
	@echo "-----------------------------------"
	@echo "Create Namespace and install Keptn-lifecycle-toolkit"
	@echo "-----------------------------------"
	helm repo add klt https://charts.lifecycle.keptn.sh
	helm repo update
	helm upgrade --install keptn klt/klt -n $(TOOLKIT_NAMESPACE) --create-namespace --wait
	kubectl apply -f support/keptn/keptnconfig.yaml -n $(TOOLKIT_NAMESPACE)

.PHONY: install-observability
install-observability:
	kubectl create namespace $(TOOLKIT_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	make -C support/observability install

.PHONY: port-forward-jaeger
port-forward-jaeger:
	make -C support/observability port-forward-jaeger

.PHONY: port-forward-grafana
port-forward-grafana:
	make -C support/observability port-forward-grafana GRAFANA_PORT_FORWARD=$(GRAFANA_PORT_FORWARD)

.PHONY: uninstall-observability
uninstall-observability:
	make -C support/observability uninstall

.PHONY: uninstall
uninstall: uninstall-observability
	@echo "-----------------------------------"
	@echo "Uninstall Keptn-lifecycle-toolkit"
	@echo "-----------------------------------"
	helm uninstall keptn -n $(TOOLKIT_NAMESPACE)
	kubectl delete ns $(TOOLKIT_NAMESPACE) --ignore-not-found=true

.PHONY: restart-lifecycle-toolkit
restart-lifecycle-toolkit:
	@echo ""
	@echo "----------------------------------"
	@echo "Restart Keptn Lifecycle Controller"
	@echo "----------------------------------"
	kubectl rollout restart deployment -n "$(TOOLKIT_NAMESPACE)" -l control-plane=lifecycle-operator
	kubectl rollout status deployment -n "$(TOOLKIT_NAMESPACE)" -l control-plane=lifecycle-operator --watch
	kubectl rollout restart deployment -n "$(TOOLKIT_NAMESPACE)" -l component=scheduler
	kubectl rollout status deployment -n "$(TOOLKIT_NAMESPACE)" -l component=scheduler --watch


kubectl port-forward --namespace makaba-reposter-staging svc/postgres 5433:5432
psql -h 127.0.0.1 -U postgres -d postgres -p 5433
kubectl exec svc/postgres -ti -- sh
kubectl exec svc/postgres -ti -n makaba-reposter-staging -- psql -U postgres
kubectl exec svc/postgres -ti -n makaba-reposter-prod -- psql -U postgres

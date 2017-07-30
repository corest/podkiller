#1/bin/bash

while true; 
do
	echo "-------------------------------------------------------"
       	kubectl get pods --all-namespaces | awk '{print "namespace: " $1 "  |  podname: " $2}'  
	sleep 2
done

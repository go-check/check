
$(.DEFAULT_GOAL) $(MAKECMDGOALS): subdirs

subdirs:
	make -C src $(MAKECMDGOALS)
	make -C src/local $(MAKECMDGOALS)

.PHONY: subdirs

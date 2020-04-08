kubedev:
	go build .

clean: 
	rm kubedev

install:
	mkdir -p $(DESTDIR)/usr/bin
	install -m 0755 kubedev $(DESTDIR)/usr/bin/kubedev
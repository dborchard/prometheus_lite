### Prometheus Lite

#### History
> Prometheus was developed at SoundCloud starting in 2012, when the company discovered that its existing metrics
> and monitoring tools. Specifically, they identified needs that Prometheus was built to meet including: a
> multi-dimensional data model, operational simplicity, scalable data collection, and a powerful query language, all in a
> single tool. The project was open-source from the beginning and began to be used by Boxever and Docker users as well,
> despite not being explicitly announced. Prometheus was inspired by the monitoring tool Borgmon used at Google

#### Code Quality
My overall opinion:

- Code package names are good
- Pointers of an Object Field is returned in some Functions. This is not a good practice. Feels kind of like C.
- Parser is simple
- PromQL is simple
- Need to revisit histogram and vectors.
- Need to revisit TSDB
- Need to revisit Fanout storage and distributed aspect of it.

#### Core Features

- [ ] Custom Parse (good starting point for tiny parser)
- [x] PromQL using AST
- [ ] Storage Adapter Layer + Fanout
- [ ] TSDB
- [ ] Rules Manager and Notifier
- [ ] Web UI + REST API

#### How to run

Right now, we have a basic test in [here](pkg/a_web/api/v1/api_test.go)

It computes 6/3 and returns 2.
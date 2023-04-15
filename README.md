# Parcour

Parcour is a Go module with goal of making parallelism and concurrency patterns easy.

The name is a mashup of parallelism and concurrency: **par**_allelism_&nbsp;**co**_nc_**ur**_rency_, and the original/french spelling of [parkour](https://en.wikipedia.org/wiki/Parkour).


## JobGroup for Structured Concurrency

The `jobgroup.JobGroup` type is similar to a nursery as described in [notes on structured concurrency](https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/).

From a Go perspective, it is meant to be passed around instead of a context, and functions use the `Go` method instead of the `go` keyword.  
This expands user driven cancellation and timeouts, benefits of a context, to also allow user driven concurrency limits, and ensuring goroutines have
an explicit scope/lifetime.

## Zync

The zync module exposes basic synchronization primitives based on their counterparts from the Go standard sync module
that are less error prone and use generics.  
See each type's documentation for details.

## Stability

In accordance with semantic versioning, no stability guarantees are made for `0.x` releases.  
Patch releases will try to keep things compatible as much as possible, but break things if required
to fix bugs.

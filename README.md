# Parcour

Parcour is a Go module with goal of making parallelism and concurrency patterns easy.

The name, which is the original/french spelling of [parkour](https://en.wikipedia.org/wiki/Parkour),
is a mashup of parallelism and concurrency: **par**_allelism_&nbsp;**co**_nc_**ur**_rency_.

## Zync

The zync module exposes basic synchronization primitives based on their counterparts from the Go standard sync module
that are less error prone and use generics.  
See each type's documentation for details.

## Stability

In accordance with semantic versioning, no stability guarantees are made for `0.x` releases.  
Patch releases will try to keep things compatible as much as possible, but break things is required
to fix bugs.

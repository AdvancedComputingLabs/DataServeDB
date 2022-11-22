


> **NOTE:** This document is provisional at the moment and subject to changes.

<!---
<table >
    <tr>
        <td><b>Creator and Author:</b></td>
        <td>Habib Yousuf AY</td>
    </tr>
    <tr>
        <td><b>Version:</b></td>
        <td>1.0</td>
    </tr>
    <tr>
        <td><b>Status:</b></td>
        <td>Draft</td>
    </tr>
</table>
--->

# DBRestQL

<!-- toc -->

- [Introduction](#introduction)
    - [Current Limitations](#current-limitations)
- [Concepts](#concepts)

<!-- /toc -->

## Introduction

> **TODO:** Introduction. Description. Why create DBRestQL and not use OData or GraphQL?

### Current Limitations

## Main Concepts

### Namespaces

Database is divided in namespaces for each database structure object type, for example: `db_name/tables` or `db_name/files`. Tables is the namespace for tables, files is the namespace for directories and files in this example.

Namespaces cannot be created or deleted by users. They are part of the database structure. A database should list namespaces it supports or are in use.

### Standard RESTful API

DBRestQL standardizes the RESTful API for database operations. It is a specification for how a client should request that resources be fetched or modified, and how a server should respond to those requests.

#### DBRestQL RESTful API for individual namespaces are as follows:
- [`tables`](dbrestql-tables.md)
- [`files`](dbrestql-files.md)
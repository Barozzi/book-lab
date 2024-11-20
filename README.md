# Book-Lab

A project to evaluate using Golang to build an API adapter for an external API.

This service wraps the Google Book API and provides a simplified interface for 
a frontend service to query book data. It is built in such a way that the external
service could be swapped out for another Book API such as Amazon.

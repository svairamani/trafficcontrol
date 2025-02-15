..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

.. _to-api-jobs:

********
``jobs``
********

``GET``
=======
Retrieve content invalidation jobs.

:Auth. Required: Yes
:Roles Required: None\ [#tenancy]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------+
	| Name            | Required | Description                                                                                                          |
	+=================+==========+======================================================================================================================+
	| assetUrl        | no       | Return only invalidation jobs that operate on URLs by matching this regular expression                               |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------+
	| createdBy       | no       | Return only invalidation jobs that were created by the user with this username                                       |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------+
	| deliveryService | no       | Return only invalidation jobs that operate on the :term:`Delivery Service` with this :ref:`ds-xmlid`                 |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------+
	| dsId            | no       | Return only invalidation jobs pending on the :term:`Delivery Service` identified by this integral, unique identifier |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------+
	| id              | no       | Return only the single invalidation job identified by this integral, unique identifer                                |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------+
	| keyword         | no       | Return only invalidation jobs that have this "keyword" - only "PURGE" should exist                                   |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------+
	| userId          | no       | Return only invalidation jobs created by the user identified by this integral, unique identifier                     |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------+


.. code-block:: http
	:caption: Request Example

	GET /api/1.4/jobs?id=3&dsId=1&userId=2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: python-requests/2.20.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:assetUrl:        A regular expression - matching URLs will be operated upon according to ``keyword``
:createdBy:       The username of the user who initiated the job
:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` on which this job operates
:id:              An integral, unique identifier for this job
:keyword:         A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetUrl`` until it is removed (or its Time to Live expires)

:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began, in a non-standard format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: gH41oEi2zrd3y8yo+wfohn4/oHU098RpyPnqBzU7HlLUDkMOPKjAZnamcYqfdy7yDCFDUcgqkvbFAvnljxyb8w==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 18 Jun 2019 19:47:30 GMT
	Content-Length: 186

	{ "response": [{
		"assetUrl": "http://origin.infra.ciab.test/.*",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 3,
		"keyword": "PURGE",
		"parameters": "TTL:2h",
		"startTime": "2019-06-18 21:28:31+00"
	}]}


``POST``
========
.. versionadded:: 1.4

Creates a new content invalidation job.

.. caution:: Creating a content invalidation job immediately triggers a CDN-wide revalidation update. In the case that the global :term:`Parameter` ``use_reval_pending`` has a value of exactly ``"0"``, this will instead trigger a CDN-wide "Queue Updates". This means that content invalidation jobs become active **immediately** at their ``startTime`` - unlike most other configuration changes they do not wait for a :term:`Snapshot` or a "Queue Updates". Furthermore, if the global :term:`Parameter` ``use_reval_pending`` *is* ``"0"``, this will cause all pending configuration changes to propagate to all :term:`cache servers` in the CDN. Take care when using this endpoint.

:Auth. Required: Yes
:Roles Required: "operations" or "admin"\ [#tenancy]_
:Response Type:  Object

Request Structure
-----------------
:deliveryService: This should either be the integral, unique identifier of a :term:`Delivery Service`, or a string containing an :ref:`ds-xmlid`
:startTime: This can be a string in the legacy ``YYYY-MM-DD HH:MM:SS`` format, or a string in :rfc:`3339` format, or a string representing a date in the same non-standard format as the ``last_updated`` fields common in other API responses, or finally it can be a number indicating the number of milliseconds since the Unix Epoch (January 1, 1970 UTC). This date must be in the future, but unlike :ref:`to-api-user-current-jobs` (or the PUT_ method of this endpoint), it is not required to be within two days from the time of creation.
:regex: A regular expression that will be used to match the path part of URIs for content stored on :term:`cache servers` that service traffic for the :term:`Delivery Service` identified by ``deliveryService``.
:ttl: Either the number of hours for which the content invalidation job should remain active, or a "duration" string, which is a sequence of numbers followed by units. The accepted units are:

	- ``h`` gives a duration in hours
	- ``m`` gives a duration in minutes
	- ``s`` gives a duration in seconds
	- ``ms`` gives a duration in milliseconds
	- ``us`` (or ``µs``) gives a duration in microseconds
	- ``ns`` gives a duration in nanoseconds

	These durations can be combined e.g. ``2h45m`` specifies a TTL of two hours and forty-five minutes - however note that durations are always rounded up to the nearest hour so that e.g. ``121m`` becomes three hours. TTLs cannot ever be negative, obviously.

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/jobs HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: python-requests/2.20.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 80
	Content-Type: application/json

	{
		"deliveryService": "demo1",
		"startTime": 1560893311219,
		"regex": "/.*",
		"ttl": "121m"
	}

Response Structure
------------------
:assetUrl:        A regular expression - matching URLs will be operated upon according to ``keyword``
:createdBy:       The username of the user who initiated the job
:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` on which this job operates
:id:              An integral, unique identifier for this job
:keyword:         A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetUrl`` until it is removed (or its Time to Live expires)

:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began, in a non-standard format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Location: https://trafficops.infra.ciab.test/api/1.4/jobs?id=3
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: nB2xg2IqO56rLT8dI4+KZgxOsTe5ShctG1U8epRsY9NyyMIpx8TZYt5MrO2QikuYh+NnyoR6V0VICCnGCKZpKw==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 18 Jun 2019 19:37:06 GMT
	Content-Length: 238

	{
		"alerts": [
			{
				"text": "Invalidation Job creation was successful",
				"level": "success"
			}
		],
		"response": {
			"assetUrl": "http://origin.infra.ciab.test/.*",
			"createdBy": "admin",
			"deliveryService": "demo1",
			"id": 3,
			"keyword": "PURGE",
			"parameters": "TTL:2h",
			"startTime": "2019-06-18 21:28:31+00"
		}
	}


``PUT``
=======
.. versionadded:: 1.4

Replaces an existing content invalidation job with a new one provided in the request. This method of editing a content invalidation job does not prevent the requesting user from changing fields that normally only have one value. Use with care.

.. caution:: Modifying a content invalidation job immediately triggers a CDN-wide revalidation update. In the case that the global :term:`Parameter` ``use_reval_pending`` has a value of exactly ``"0"``, this will instead trigger a CDN-wide "Queue Updates". This means that content invalidation jobs become active **immediately** at their ``startTime`` - unlike most other configuration changes they do not wait for a :term:`Snapshot` or a "Queue Updates". Furthermore, if the global :term:`Parameter` ``use_reval_pending`` *is* ``"0"``, this will cause all pending configuration changes to propagate to all :term:`cache servers` in the CDN. Take care when using this endpoint.

:Auth. Required: Yes
:Roles Required: "operations" or "admin"\ [#tenancy]_
:Response Type:  Object

Request Structure
-----------------
.. table:: Query Parameters

	+------+----------+--------------------------------------------------------------------------------+
	| Name | Required | Description                                                                    |
	+======+==========+================================================================================+
	| id   | yes      | The integral, unique identifier of the content invalidation job being modified |
	+------+----------+--------------------------------------------------------------------------------+

:assetUrl: A regular expression - matching URLs will be operated upon according to ``keyword``

	.. note:: Unlike in the payloads of POST_ requests to this endpoint, this must be a **full** URL regular expression, as it is **not** combined with the :ref:`ds-origin-url` of the :term:`Delivery Service` identified by ``deliveryService``.

:createdBy:       The username of the user who initiated the job\ [#readonly]_
:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` on which this job operates\ [#readonly]_ - unlike POST_ request payloads, this cannot be an integral, unique identifier
:id:              An integral, unique identifier for this job\ [#readonly]_
:keyword:         A keyword that represents the operation being performed by the job. It can have any (string) value, but the only value with any meaning to Traffic Control is:

	PURGE
		This job will prevent caching of URLs matching the ``assetUrl`` until it is removed (or its Time to Live expires)

:parameters: A string containing space-separated key/value pairs - delimited by colons (:kbd:`:`\ s) representing parameters associated with the job. In practice, any string can be passed as a job's ``parameters``, but the only value with meaning is a single key/value pair indicated a :abbr:`TTL (Time To Live)` in hours in the format :file:`TTL:{hours}h`, and any other type of value may cause components of Traffic Control to work improperly or not at all.
:startTime:  This can be a string in the legacy ``YYYY-MM-DD HH:MM:SS`` format, or a string in :rfc:`3339` format, or a string representing a date in the same non-standard format as the ``last_updated`` fields common in other API responses, or finally it can be a number indicating the number of milliseconds since the Unix Epoch (January 1, 1970 UTC). This **must** be in the future, but only by no more than two days.

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/jobs?id=3 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: python-requests/2.20.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 188
	Content-Type: application/json

	{
		"assetUrl": "http://origin.infra.ciab.test/.*",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 3,
		"keyword": "PURGE",
		"parameters": "TTL:360h",
		"startTime": "2019-06-20 18:33:40+00"
	}

Response Structure
------------------
:assetUrl:        A regular expression - matching URLs will be operated upon according to ``keyword``
:createdBy:       The username of the user who initiated the job
:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` on which this job operates
:id:              An integral, unique identifier for this job
:keyword:         A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetUrl`` until it is removed (or its Time to Live expires)

:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began, in a non-standard format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: +P1PTav4ZBoiQcCqQnUqf+J0dCfQgVj8mzzKtUCA69mWYulya9Bjf6BUd8Aro2apmpgPBkCEA5sITJV1tMYA0Q==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Jun 2019 13:38:59 GMT
	Content-Length: 234

	{ "alerts": [{
		"text": "Content invalidation job updated",
		"level": "success"
	}],
	"response": {
		"assetUrl": "http://origin.infra.ciab.test/.*",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 3,
		"keyword": "PURGE",
		"parameters": "TTL:360h",
		"startTime": "2019-06-20 18:33:40+00"
	}}


``DELETE``
==========
.. versionadded:: 1.4

Deletes a content invalidation job.

.. tip:: Content invalidation jobs that have passed their :abbr:`TTL (Time To Live)` are not automatically deleted - for record-keeping purposes - so use this to clean up old jobs that are no longer useful.

.. caution:: Deleting a content invalidation job immediately triggers a CDN-wide revalidation update. In the case that the global :term:`Parameter` ``use_reval_pending`` has a value of exactly ``"0"``, this will instead trigger a CDN-wide "Queue Updates". This means that content invalidation jobs become active **immediately** at their ``startTime`` - unlike most other configuration changes they do not wait for a :term:`Snapshot` or a "Queue Updates". Furthermore, if the global :term:`Parameter` ``use_reval_pending`` *is* ``"0"``, this will cause all pending configuration changes to propagate to all :term:`cache servers` in the CDN. Take care when using this endpoint.

:Auth. Required: Yes
:Roles Required: "operations" or "admin"\ [#tenancy]_
:Response Type:  Object

Request Structure
-----------------
.. table:: Query Parameters

	+------+----------+--------------------------------------------------------------------------------+
	| Name | Required | Description                                                                    |
	+======+==========+================================================================================+
	| id   | yes      | The integral, unique identifier of the content invalidation job being modified |
	+------+----------+--------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/jobs?id=3 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: python-requests/2.20.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0

Response Structure
------------------
:assetUrl:        A regular expression - matching URLs will be operated upon according to ``keyword``
:createdBy:       The username of the user who initiated the job
:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` on which this job operates
:id:              An integral, unique identifier for this job
:keyword:         A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetUrl`` until it is removed (or its Time to Live expires)

:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began, in a non-standard format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: FqfziXJYYwHb84Fac9+p4NEY3EsklYxe94wg/VOmlXk4R6l4SaPSh015CChPt/yT72MsWSETnIuRD9KtoK4I+w==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 18 Jun 2019 22:55:15 GMT
	Content-Length: 234

	{ "alerts": [
		{
			"text": "Content invalidation job was deleted",
			"level": "success"
		}
	],
	"response": {
		"assetUrl": "http://origin.infra.ciab.test/.*",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 3,
		"keyword": "PURGE",
		"parameters": "TTL:36h",
		"startTime": "2019-06-20 18:33:40+00"
	}}


.. [#tenancy] When viewing content invalidation jobs, only those jobs that operate on a :term:`Delivery Service` visible to the requesting user's :term:`Tenant` will be returned. Likewise, creating a new content invalidation job requires that the target :term:`Delivery Service` is modifiable by the requesting user's :term:`Tenant`. However, when modifying or deleting an existing content invalidation job, the operation can be completed if and only if the requesting user's :term:`Tenant` is the same as the job's :term:`Delivery Service`'s :term:`Tenant` or a descendant thereof, **and** if the requesting user's :term:`Tenant` is the same as the :term:`Tenant` of the *user who initially created the job* or a descendant thereof.
.. [#readonly] This field must exist, but it must *not* be different than the same field of the existing job (i.e. as seen in a GET_ response)

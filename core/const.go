/*
Copyright © 2022 The K9CLI Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

const FIRST_YEAR = 2021

const (
	DB_INDEX_POSITION_CUSTOMERID = 1
	DB_INDEX_POSITION_ACCOUNT    = 4
	DB_INDEX_POSITION_YEAR       = 5
	DB_INDEX_POSITION_MONTH      = 6
	DB_INDEX_POSITION_FILE       = 7
)

const (
	ACCESS_CAPABILITY_RESOURCE_ADMIN = `administer-resource`
	ACCESS_CAPABILITY_DELETE_DATA    = `delete-data`
	ACCESS_CAPABILITY_READ_CONFIG    = `read-config`
	ACCESS_CAPABILITY_READ_DATA      = `read-data`
	ACCESS_CAPABILITY_WRITE_DATA     = `write-data`
)

# Chef API testing

For each class of chef api object class (organizations, users, cookbooks, etc and each rest end point)
* Documentation updated
* Exercise the defined functions - mostly crud
* Define input structs for the calls to the api
* Create an example of the code to call each function
* Define the expected output in a struct
* Verify the output structs are populated as expected

|Name| New | Doc update | Function tests | input structs | out struct |
|----|-----|------------|----------------|---------------|------------|
|GLOBAL LEVEL|
|authenticate_user |x|x||||
|license |x|x||||
|organizations|x|x||||
|status |x|x||||
|users|x|x||||
|ORGANIZATION LEVEL|||||||
|association_request |x|x||||
|clients|x|x||||
|containers ||||||
|cookbook|x|||||
|cookbook_download||||||
|databag|x|||||
|environments|x|||||
|groups|x|x||||
|nodes|x|||||
|policyx|x|||||
|policy_group |x|||||
|principal|U|x|x|x|x|
|roles|x|x||||
|run_list||||||
|sandboxes||x||||
|search|x|||||
|universe |x|x||||
|updated_since |x|||||
|users(in associations)|x|||||
|organization/users||||||
|acl||||||

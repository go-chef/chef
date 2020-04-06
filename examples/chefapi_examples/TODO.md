* Fix the circleci pipeline

* live testing
* sandbox

Missing api stuff
* user authentication
* display the user external_authentication_uid value

The search doc needs to be clearer. It refers to the knife search string but isn't clear that what it means

Sandbox upload gets a 403. https://stackoverflow.com/questions/28807834/chef-upload-to-sandbox-fails-with-403-unauthorized
It's not using the client, going to a returned url and give an auth code.  Not sure what to do with it.
embedded/lib/ruby/gems/2.6.0/gems/chef-15.6.10/lib/chef/cookbook_uploader.rb has the ruby code to do this stuff

2879448ee0bb891f2103f21891234b2f:{Url:https://localhost:443/bookshelf/organization-87e84b155d98af0a68fb7e69e7adee0a/checksum-2879448ee0bb891f2103f21891234b2f?AWSAccessKeyId=3d54e25d43d14490faa55cf85826a2732e3c9411&Expires=1577447683&Signature=cQYAocsq%2B15ztGX3zEcmJ1C9HKk%3D Upload:true} 2b26f685cf7bce1b634787df20607aeb:{Url:https://localhost:443/bookshelf/organization-87e84b155d98af0a68fb7e69e7adee0a/checksum-2b26f685cf7bce1b634787df20607aeb?AWSAccessKeyId=3d54e25d43d14490faa55cf85826a2732e3c9411&Expires=1577447683&Signature=HBsn1DnQWuItPVxQNtCG7Rjc18w%3D Upload:true} 73008070bb42cbbe05dfa5427a8416b1:{Url:https://localhost:443/bookshelf/organization-87e84b155d98af0a68fb7e69e7adee0a/checksum-73008070bb42cbbe05dfa5427a8416b1?AWSAccessKeyId=3d54e25d43d14490faa55cf85826a2732e3c9411&Expires=1577447683&Signature=5RrhZNMvUU%2B4D5fZkYsuYdClWdQ%3D Upload:true}]}

Uploading: 2879448ee0bb891f2103f21891234b2f --->  {https://localhost:443/bookshelf/organization-87e84b155d98af0a68fb7e69e7adee0a/checksum-2879448ee0bb891f2103f21891234b2f?AWSAccessKeyId=3d54e25d43d14490faa55cf85826a2732e3c9411&Expires=1577447683&Signature=cQYAocsq%2B15ztGX3zEcmJ1C9HKk%3D true}

after request &{PUT https://localhost:443/bookshelf/organization-87e84b155d98af0a68fb7e69e7adee0a/checksum-2879448ee0bb891f2103f21891234b2f?AWSAccessKeyId=3d54e25d43d14490faa55cf85826a2732e3c9411&Expires=1577447683&Signature=cQYAocsq%2B15ztGX3zEcmJ1C9HKk%3D HTTP/1.1 1 1 map[X-Ops-Authorization-6:[bDFCEcSiDHAqCqwtvwNP8rX/0HgDwO4cCvaNJQT7uw==] Content-Type:[application/octet-stream] Method:[PUT] X-Ops-Userid:[pivotal] X-Ops-Authorization-2:[qMNYPwz2Ym6awXiCgOy6YXxGZeCvfvTngGPvdDa/FPVJLt5MUdYBNOBk/L1W] X-Ops-Authorization-5:[OG9NI2aOPOnah512Cav2l7H/mgdM4gr3+jyuz/sGQVxAv5jriu086zKqU/EJ] X-Ops-Content-Hash:[RWKT0i3TZ5n3l3km2zgS6JlPs9M=] X-Ops-Timestamp:[2019-12-27T11:39:43Z] X-Ops-Authorization-4:[PY95MHZVNgaH3gcYj+dD3GHCtDlk44C4vHEHv3N5OVTfrcs7Mi3tuprPARuh] Accept:[application/json] X-Ops-Authorization-1:[wFXZeEpvMefVCPDVGUZL/zNMlqx94yN5L78rq13TJticQNXCXJg45Y74fwKY] X-Ops-Authorization-3:[7V1BquFUsAJxsKkhiGsvT0l6qUTNFIKOGc2KOH6p0cGpPlzaDBoEyd8Qer0U] X-Chef-Version:[11.12.0] X-Ops-Sign:[algorithm=sha1;version=1.0]] {0xc42006ddd0} 0x604110 128 [] false localhost:443 map[] map[] <nil> map[]   <nil> <nil> <nil> <nil>}
before upload
after upload PUT https://localhost:443/bookshelf/organization-87e84b155d98af0a68fb7e69e7adee0a/checksum-2879448ee0bb891f2103f21891234b2f?AWSAccessKeyId=3d54e25d43d14490faa55cf85826a2732e3c9411&Expires=1577447683&Signature=cQYAocsq%2B15ztGX3zEcmJ1C9HKk%3D: 403

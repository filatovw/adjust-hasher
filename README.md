Hasher
======

Here is a [challenge](./docs/challenge.md)


Build application and put into `./hasher/bin/`

    make build

How to use
==========

Execute against one address:

    ./hasher/bin/hasher adjust.com

    http://adjust.com 31d8d84f3aa297e1e495d9c5b2942076

Execute against multiple addresses (default concurrency: 10):

    ./hasher/bin/hasher adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com

    http://adjust.com 31d8d84f3aa297e1e495d9c5b2942076
    http://twitter.com 7a7b6b59a6b78c7be159321cd06bbee5
    http://yandex.com 646668d7eedd662f4d7699fcad757ff4
    http://baroquemusiclibrary.com b93c46c58a7b9c466cb983a6afd14fd8
    http://facebook.com 73cd843f772c59501e692d8b7a2fbb8a
    http://reddit.com/r/notfunny e8e40dcbccc41e733f1e93f73c3d9967
    http://reddit.com/r/funny d79ab62de3b7bd7b0f3df4d3210c3073
    http://yahoo.com 4fee9e8ef17b714de5551ac36594d98c
    http://google.com dffc2c62185c95f131f34b4889224a02

Execute with concurrency parameter:

    ./hasher/bin/hasher -parallel 200  adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com

    http://adjust.com 31d8d84f3aa297e1e495d9c5b2942076
    http://google.com a426e33cbafb5a2805ef6c196f8bc8e8
    http://twitter.com c7335cd92288b1faf58ad9c9263d5552
    http://yandex.com 867cde3e70786afcc8952345b04ada87
    http://facebook.com 43ce9850da9c849ce485966d371a89bf
    http://reddit.com/r/notfunny f9d5d8c2d5cdbda5fd65bc6615c99632
    http://baroquemusiclibrary.com 271d149e560fe5303649b0fa1e1768ad
    http://yahoo.com f4dbd94ecee45bde7e07262dc0ef0178
    http://reddit.com/r/funny 5d2d77be0371fed9e4f5675da50da883

Execute in debug mode:

    ./hasher/bin/hasher -debug -parallel 2  adjust.com google.com facebook.com

    hasher ### 2019/02/06 18:08:34.131314 /Users/vadimfilatov/projects/experiments/go/src/github.com/filatovw/adjust-hasher/hasher/app/pool.go:38: worker 2: http://adjust.com :: 31d8d84f3aa297e1e495d9c5b2942076
    http://adjust.com 31d8d84f3aa297e1e495d9c5b2942076
    hasher ### 2019/02/06 18:08:34.192738 /Users/vadimfilatov/projects/experiments/go/src/github.com/filatovw/adjust-hasher/hasher/app/pool.go:38: worker 1: http://google.com :: 90a40c66e85cbcfc455f83b35a55272e
    http://google.com 90a40c66e85cbcfc455f83b35a55272e
    hasher ### 2019/02/06 18:08:34.939621 /Users/vadimfilatov/projects/experiments/go/src/github.com/filatovw/adjust-hasher/hasher/app/pool.go:38: worker 2: http://facebook.com :: 7ff5559bf9fb48e7cbbf4daf7d336d56
    http://facebook.com 7ff5559bf9fb48e7cbbf4daf7d336d56

If request failed write error to output:

    ./hasher/bin/hasher -parallel 2  adjust.com google.com facebook.c

    http://adjust.com 31d8d84f3aa297e1e495d9c5b2942076
    error: download: Get http://facebook.c: dial tcp: lookup facebook.c: no such host
    http://google.com df035843c0a1aa01b44ef4c61161bf27


cur=$(pwd)
cmd='go mod tidy && go get -u && go mod tidy'
# cmd='go fix ./...'
cmd='go mod tidy'
# --diff: diff文件
# cmd='go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix ./...'

# go clean -modcache
find . -name "go.mod" -type f | while read -r mod_file; do
    dir=$(dirname "$mod_file")
    if [ "$dir" != "." ]; then
        # go work use "$dir" 2>/dev/null || echo "  跳过"
        cd $dir && pwd && eval $cmd;
        cd $cur
    fi
done

# ids="276 279 294 295 300 303 304 306 307 308 311 312 314 315 316 317 319 320 324 325 326 339 340 341 344 348 351 352 353 354 355 382 383 386 389 391 392 393 394 395 396 397 398 399 402 403 405 406 407 409 410 412 413 414 415 442 443 444 445 446"
# for id in $ids; do
#     curl --request GET --url "http://localhost:58090/weibo/rss/clear?id=${id}"
# done
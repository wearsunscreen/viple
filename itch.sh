rm -rf ~/viple-zip
cp ~/viple/ ~/viple-zip
rm ~/viple-zip/*.go
rm ~/viple-zip/viple
rm ~/viple-zip/README.md
rm ~/viple-zip/*.sh
rm ~/viple-zip/.gitignore
zip ~/viple.zip ~/viple-zip/*.* 

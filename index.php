<?php

//The name of the directory that we need to create.
$directoryName = '.cache';

//Check if the directory already exists.
if(!is_dir($directoryName)){
    //Directory does not exist, so lets create it.
    mkdir($directoryName, 7777);
    echo 'Folder Created';
} else {
    echo 'Folder exist';
}

//The name of the directory that we need to create.
$directoryName = '.cache/googleplay';

//Check if the directory already exists.
if(!is_dir($directoryName)){
    //Directory does not exist, so lets create it.
    mkdir($directoryName, 7777);
    echo 'Folder Created';
} else {
    echo 'Folder exist';
}

$myfile = fopen(".cache/googleplay/token.json", "w") or die("Unable to open file!");
$txt = "";
fwrite($myfile, $txt);
fclose($myfile);

?>

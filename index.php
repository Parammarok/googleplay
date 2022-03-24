<?php

//The name of the directory that we need to create.
$directoryName = '.cache';

//Check if the directory already exists.
if(!is_dir($directoryName)){
    //Directory does not exist, so lets create it.
    mkdir($directoryName, 0777);
    echo 'Folder Created';
} else {
    echo 'Folder exist';
}


?>

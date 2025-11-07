INSERT INTO admin (
	admin_email
) VALUES (
    $1
) 
RETURNING *;

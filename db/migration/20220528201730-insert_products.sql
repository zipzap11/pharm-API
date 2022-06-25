
-- +migrate Up
INSERT INTO products (id, category_id, name, price, weight, description, image_url) VALUES
(
    1, 
    3, 
    'Rhinos Junior Sirup 60 ml', 
    68700, 
    60, 
    'RHINOS JUNIOR merupakan obat dengan kandungan Pseudoephedrine HCl dan Chlorpheniramine Maleate yang digunakan untuk 
    meringankan gejala flu seperti bersin-bersin dan hidung tersumbat karena pilek. Indikasi Umum
    Obat ini digunakan untuk meringankan gejala bersin-bersin dan hidung tersumbat karena pilek. Dapat 
    diberikan bersama atau tanpa makanan.',
    'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/products%2F1.png?alt=media&token=5de3d116-9cd4-4094-bdfa-f810b812eb4b'
),
(
    2,
    3,
    'Acetylcysteine 200 mg 10 Kapsul',
    24700,
    0.2,
    'ACETYLCYSTEINE 200 MG KAPSUL adalah obat generik yang digunakan sebagai mukolitik (pengencer dahak) 
    dan antidot pada pasien yang overdosis Paracetamol. Fungsi mukolitik dari Acetylcysteine bekerja dengan 
    cara melalui gugus sulfhidril bebasnya yang membuka ikatan disulfida dalam mukoprotein, sehingga menurunkan 
    viskositas lendir/ mukus (dahak) yang membuat mukus tersebut encer sehingga mukus lebih mudah dikeluarkan. 
    Dalam penggunaan obat ini HARUS SESUAI DENGAN PETUNJUK DOKTER.',
    'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/products%2F2.png?alt=media&token=da519c3b-fbd5-427f-b9db-83a70c84107b'
),
(
    3,
    3,
    'Rhinos Neo Drops 10 ml',
    80200,
    10,
    'RHINOS NEO DROPS merupakan obat dengan kandungan Pseudoephedrine HCl. Obat ini dapat digunakan untuk 
    meringankan gejala hidung tersumbat krn flu.',
    'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/products%2F3.png?alt=media&token=01a299ca-4f8e-4ba7-b3f5-b2eb547a2fa2'
),
(
    4,
    1,
    'Ventolin Inhaler 100 mcg 200 Doses',
    175700,
    0.0001,
    'VENTOLIN INHALER merupakan obat dengan kandungan Salbutamol yang digunakan untuk mengobati penyakit pada 
    saluran pernafasan seperti asma dan penyakit paru obstruktif kronik (PPOK). Obat ini bekerja dengan cara merangsang secara selektif reseptor 
    beta-2 adrenergik terutama pada otot bronkus. hal ini menyebabkan terjadinya bronkodilatasi karena otot bronkus mengalami relaksasi. 
    Dalam penggunaan obat ini harus SESUAI DENGAN PETUNJUK DOKTER.',
    'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/products%2F4.png?alt=media&token=9bbafc5f-137e-46fc-a7e6-b290fd9bec29'
),
(
    5,
    1,
    'Lasal Expectorant Sirup 100 ml',
    58200,
    100,
    'LASAL EXPECTORANT SIRUP adalah obat yang mengandung kombinasi Salbutamol dan Guaifenesin. Obat ini digunakan 
    untuk melegakan pernafasan dan membantu pengeluaran dahak pada kondisi asma yang disertai batuk. 
    Senyawa salbutamol merupakan obat golongan bronkodilator yang bekerja dengan cara merangsang secara selektif 
    reseptor beta-2 adrenergik terutama pada otot bronkus.',
    'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/products%2F5.png?alt=media&token=ceaf659d-5517-4783-8463-2a566db36545'
),
(
    6,
    5,
    'Diatomix Pembasmi Kutu Pada Hewan Peliharaan & Ternak 25 g',
    27500,
    25,
    'DIATOMIX PEMBASMI KUTU PADA HEWAN PELIHARAAN & TERNAK merupakan obat kutu khusus untuk hewan peliharaan dan 
    ternak yang mengandung Amorphous Activated Diatomaceous Earth (ADE)',
    'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/products%2F6.png?alt=media&token=b4e97ea7-3232-4419-aec6-313f4875c6cb'
),
(
    7,
    5,
    'Beauty Premium Cat Food Formulated for Adult Cats & Kittens 20 kg',
    491000,
    20000,
    'PRODUK INI HANYA DIKONSUMSI UNTUK HEWAN. Makanan untuk kucing dewasa dan kitten atau anakan kucing',
    'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/products%2F7.png?alt=media&token=bfdcc724-36cd-4726-95ba-88f46c04c9c0'
);
-- +migrate Down
DELETE FROM products WHERE id in (1, 2, 3, 4, 5, 6, 7);
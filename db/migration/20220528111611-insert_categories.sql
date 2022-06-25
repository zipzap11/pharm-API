
-- +migrate Up
INSERT INTO CATEGORIES (id, name, image_url) VALUES
    (1, 'asma', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/asma.png?alt=media&token=be73404c-cb58-4c97-8bb5-e43c7487961a'),
    (2, 'jantung', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/jantung.png?alt=media&token=6f630283-7dc5-4f28-ab7c-d4cfb2d26d0d'),
    (3, 'batuk & flu', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/batuk_flu.png?alt=media&token=2c5074f9-5fb5-4d5a-af53-bd4bbb0f3262'),
    (4, 'herbal', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/herbal.png?alt=media&token=81e66fb5-85ca-4140-954f-d610c718443f'),
    (5, 'hewan', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/hewan.png?alt=media&token=6ecdd8ae-ff11-43fb-99c4-a6966037e55c'),
    (6, 'pria', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/pria.png?alt=media&token=c1a074e0-3670-4c79-8c9c-e03f01f58712'),
    (7, 'minyak', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/minyak.png?alt=media&token=c1441c7c-7374-48db-97f3-a2b514a02894'),
    (8, 'P3K', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/p3k.png?alt=media&token=8a3ac506-84d1-4887-a116-2b05384c43de'),
    (9, 'mulut & tenggorokan', 'https://firebasestorage.googleapis.com/v0/b/linkedin-ef4e1.appspot.com/o/mulut_tenggorokan.png?alt=media&token=ea850f70-8de7-420a-9013-3f5685cf1487')
;
-- +migrate Down
DELETE FROM CATEGORIES WHERE id in (1,2,3,4,5,6,7,8,9);
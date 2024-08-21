VMAF
source: https://github.com/Netflix/vmaf/blob/master/resource/doc/ffmpeg.md

vmaf  quality
---------------------
0-20   bad
20-50  poor
50-70  fair/good
70-100 good/excellent

example: ffmpeg -i sravnivaemoe.mp4 -i etalon.mp4 -lavfi libvmaf=model_path="c\\:/Users/pemaltynov/go/src/github.com/Netflix/vmaf/model/vmaf_v0.6.1.json" -f null -

read more: https://jina-liu.medium.com/a-practical-guide-for-vmaf-481b4d420d9c

------------------------------------------------------------------------------------------------------------------------

SSIM
source: https://www.testdevlab.com/blog/full-reference-quality-metrics-vmaf-psnr-and-ssim

ssim    video degradation
--------------------
0.95-0.97    low
0.97-1       minimal

psnr(dB) quality
------------------------
20-30     bad
30-33     poor
33-38     fair
38+       good/excellent

example: ffmpeg -i sravnivaemoe.mpg -i etalon.mpg -lavfi  "ssim;[0:v][1:v]psnr" -f null -
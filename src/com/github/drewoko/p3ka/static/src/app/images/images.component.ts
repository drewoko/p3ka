import {Component, Input, OnChanges, Output, EventEmitter, OnInit} from "@angular/core";
import {Image} from "./image";
import {Location} from "@angular/common";
import {Router} from "@angular/router";
import {ImageService} from "./image.service";

@Component({
    selector: 'images',
    templateUrl: './images.component.html',
    styleUrls: ['./images.component.css']
})
export class ImagesComponent implements OnInit {

    @Input('images')
    images: Image[] = [];

    @Input('big')
    big: boolean;

    @Output()
    onNewImageLoadRequested = new EventEmitter();

    selectedImage: Image = null;

    constructor(private location: Location, private router: Router, private imageService: ImageService) {
    }

    ngOnInit() {
        this.imageService.forceOpenImageAnnounced$.subscribe(image => {
            this.openImage(image);
        });
    }

    openImage(selectedImage: Image): void {
        this.selectedImage = selectedImage;
        //todo: fix it
        if(selectedImage != null) {
            this.location.go("/show/" + selectedImage.id);
        } else {
            if(this.router.url.includes("/show")) {
                this.location.go("/user/" + this.images[0].name);
            } else {
                this.location.go(this.router.url);
            }
        }
    }

    imageEvent(event: MouseEvent): void {
        if(event.toElement.nodeName == "IMG") {
            let foundNextImage: Image = this.images[this.images.indexOf(this.selectedImage) + 1];

            if(foundNextImage == null) {
                this.imageService.imageLoadRequest.next();
            } else {
                this.openImage(foundNextImage);
            }
        } else {
            this.openImage(null);
        }
    }
}
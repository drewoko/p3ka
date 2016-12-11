import {Component, Input, OnInit} from "@angular/core";
import {Image} from "./image";
import {Location} from "@angular/common";
import {Router} from "@angular/router";
import {ImageService} from "./image.service";

@Component({
    selector: 'images',
    templateUrl: './images.component.html',
    styleUrls: ['./images.component.css'],
    host: {
        '(document:keydown)': 'handleKeyboardEvent($event)'
    }
})
export class ImagesComponent implements OnInit {

    @Input('images')
    images: Image[] = [];

    @Input('big')
    big: boolean;

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
            this.nextImage();
        } else {
            this.closeImage();
        }
    }

    handleKeyboardEvent(event: KeyboardEvent) {
        if(event.key == "ArrowRight") {
            this.nextImage();
        } else if (event.key == "ArrowLeft") {
            this.prevImage();
        } else if (event.key == "Escape") {
            this.closeImage();
        }
    }

    private closeImage(): void {
        this.openImage(null);
    }

    private prevImage(): void {
        this.openImage(this.images[this.images.indexOf(this.selectedImage) - 1]);
    }

    private nextImage(): void {
        let foundNextImage: Image = this.images[this.images.indexOf(this.selectedImage) + 1];

        if(foundNextImage == null) {
            this.imageService.imageLoadRequest.next();
        } else {
            this.openImage(foundNextImage);
        }
    }
}
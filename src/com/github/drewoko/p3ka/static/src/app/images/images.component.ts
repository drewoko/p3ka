import {Component, Input, OnChanges} from "@angular/core";
import {Image} from "./image";
import {Location} from "@angular/common";
import {Router} from "@angular/router";

@Component({
    selector: 'images',
    templateUrl: './images.component.html',
    styleUrls: ['./images.component.css']
})
export class ImagesComponent implements OnChanges {

    @Input('images')
    images: Image[] = [];

    @Input('big')
    big: boolean;

    @Input('open')
    open: number;

    selectedImage: Image = null;

    constructor(private location: Location, private router: Router) {}

    ngOnChanges(changes) {
        if(this.open != null && this.images.length > 0) {
            this.findImage(this.open, this.images)
                .then(image => this.openImage(image))
                .catch(err => {
                    this.open = null;
                });
        }
    }

    openImage(selectedImage: Image): void {
        this.selectedImage = selectedImage;
        //todo: fix it
        if(selectedImage != null) {
            this.location.go("/show/" + selectedImage.id);
        } else {
            if(this.open != null) {
                this.location.go("/user/" + this.images[0].name);
                this.open = null;
            } else {
                if(this.router.url.includes("/show")) {
                    this.location.go("/user/" + this.images[0].name);
                } else {
                    this.location.go(this.router.url);
                }
            }
        }
    }

    imageEvent(event: MouseEvent): void {
        this.openImage(
            event.toElement.nodeName == "IMG" ? this.images[this.images.indexOf(this.selectedImage) + 1] : null)
    }

    private findImage(id: number, images: Image[]): Promise<Image> {
        return new Promise<Image>((resolve, reject) => {
            images.forEach(image => {
              if(image.id == id) {
                  resolve(image);
              }
            });
        });
    }
}
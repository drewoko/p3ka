import {Component, Input} from "@angular/core";
import {Image} from "./image";

@Component({
    selector: 'images',
    templateUrl: './images.component.html',
    styleUrls: ['./images.component.css']
})
export class ImagesComponent {

    @Input('images')
    images: Image[];

    @Input('big')
    big: boolean = false;

    selectedImage: Image = null;

    openImage(selectedImage: Image): void {
        this.selectedImage = selectedImage;
    }

    imageEvent(event: MouseEvent): void {
        this.openImage(
            event.toElement.nodeName == "IMG" ? this.images[this.images.indexOf(this.selectedImage) + 1] : null)
    }
}
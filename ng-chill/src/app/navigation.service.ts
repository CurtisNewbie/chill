import { Injectable } from "@angular/core";
import { Router } from "@angular/router";

@Injectable({
  providedIn: "root",
})
export class NavigationService {
  constructor(private router: Router) { }

  /** Navigate to using Router*/
  public navigateToUrl(url: string, extra?: any[]): void {
    let arr: any[] = [url];
    if (extra != null) arr = arr.concat(extra);
    this.router.navigate(arr);
  }
}


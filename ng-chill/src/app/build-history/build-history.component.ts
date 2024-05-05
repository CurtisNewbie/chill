import { Component } from '@angular/core';
import { PagingController } from '../common/paging';
import { HttpClient } from '@angular/common/http';
import { Toaster } from '../toaster.service';
import { NavigationService } from '../navigation.service';

export interface BuildHist {
  id?: number
  name?: string
  buildNo?: string
  status?: string
  ctime?: number | Date
}

@Component({
  selector: 'app-build-history',
  templateUrl: './build-history.component.html',
  styleUrl: './build-history.component.css'
})
export class BuildHistoryComponent {

  data: any[] = []
  pagingController: PagingController;

  constructor(private http: HttpClient, private toaster: Toaster, private nav: NavigationService) {
  }

  fetchList() {
    this.http.post<any>("/api/build/history/list", { paging: this.pagingController.paging })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.toaster.toast(resp.msg);
            return;
          }

          this.data = [];
          if (resp.data.payload) {
            for (let r of resp.data.payload) {
              if (r.ctime) {
                r.ctime = new Date(r.ctime);
              }
              this.data.push(r);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
        error: (err) => {
          this.toaster.toast(`Request failed: ${err}`)
        },
      });

  }

  onPagingControllerReady(pc: PagingController) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  popHistDetails(u: BuildHist) {
    this.nav.navigateToUrl("/build/history/details", [
      { buildNo: u.buildNo },
    ]);
  }
}

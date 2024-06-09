import { Component, OnInit } from '@angular/core';
import { PagingController } from '../common/paging';
import { HttpClient } from '@angular/common/http';
import { Toaster } from '../toaster.service';
import { NavigationService } from '../navigation.service';
import { ActivatedRoute } from '@angular/router';

export interface ApiListBuildHistoryRes {
  id?: number                    // build history id
  name?: string                  // build name
  buildNo?: string               // build no
  commitId?: string              // build commit id
  tag?: string                   // tag
  status?: string                // built status
  startTime?: number             // build start time
  endTime?: number               // build end time
}

@Component({
  selector: 'app-build-history',
  templateUrl: './build-history.component.html',
  styleUrl: './build-history.component.css'
})
export class BuildHistoryComponent implements OnInit {

  name: string = null
  data: ApiListBuildHistoryRes[] = []
  pagingController: PagingController;

  constructor(private http: HttpClient, private toaster: Toaster,
    private nav: NavigationService, private route: ActivatedRoute) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe((params) => {
      let n = params.get("name");
      if (n) {
        this.name = n;
      }
    });
  }

  fetchList() {
    this.http.post<any>("/api/build/history/list", { paging: this.pagingController.paging, name: this.name })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.toaster.toast(resp.msg);
            return;
          }

          this.data = [];
          if (resp.data.payload) {
            for (let r of resp.data.payload) {
              this.data.push(r);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
        error: (err) => {
          this.toaster.toast(`Request failed`)
          console.log(err)
        },
      });

  }

  onPagingControllerReady(pc: PagingController) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  popHistDetails(u: ApiListBuildHistoryRes) {
    this.nav.navigateToUrl("/build/history/details", [
      { buildNo: u.buildNo },
    ]);
  }
}

import d3 from 'd3';
import React from 'react';
import ReactFauxDOM from 'react-faux-dom';
import JQuery from 'jquery';

class Chart extends React.Component {
  static propTypes = {
    data: React.PropTypes.array,
    meta: React.PropTypes.array,
    x: React.PropTypes.array,
    title: React.PropTypes.string,
    yMax: React.PropTypes.number
  }
  render() {
    var that = this;
    var m =  120, // number of samples per layer
        data = this.props.data,
        xdata = this.props.x,
        yMax = typeof this.props.yMax !== "undefined" ? this.props.yMax : d3.max(this.props.data);

    var margin = {top: 40, right: 10, bottom: 50, left: 60},
        width = 910 - margin.left - margin.right,
        height = 180 - margin.top - margin.bottom;
    
    var x = d3.scale.ordinal()
              .domain(d3.range(m))
              .rangeRoundBands([width, 0], 0.20, 0.10);

    var y = d3.scale.linear()
              .domain([0, yMax])
              .range([height, 0]);
    var detail = (<span></span>);
    
    var timeFormat = d3.time.format("%I:%M %p");
    var datetimeFormat = d3.time.format("%d/%m/%y %I:%M %p");
    const xAxisFormatter = (t) => {
      return timeFormat(new Date(xdata[t] * 1000));
    };

    const yFormatter = this.props.yFormatter;

    var xAxis = d3.svg.axis()
                  .scale(x)
                  .tickValues(d3.range(0, xdata.length, 12))
                  .tickSize(2, 0)
                  .tickPadding(6)
                  .tickFormat(xAxisFormatter)
                  .orient("bottom");

    var yAxis = d3.svg.axis()
                  .scale(y)
                  .ticks(3)
                  .tickSize(2, 0)
                  .tickPadding(6)
                  .tickFormat(yFormatter)
                  .orient("left");
    
    const node = ReactFauxDOM.createElement('svg');
    const svg = d3.select(node)
                  .attr("width", width + margin.left + margin.right)
                  .attr("height", height + margin.top + margin.bottom)
                  .append("g")
                  .attr("transform", "translate(0," + margin.top + ")");

    var rect = svg.selectAll("rect")
                  .data(data)
                  .enter().append("rect")
                  .attr("x", function(d, i) { return x(i) + margin.left; })
                  .attr("width", x.rangeBand())
                  .attr("y", function(d) { return y(d); })
                  .attr("height", function(d) { return y(0) - y(d); })
                  .on("mouseover", (d, i) => {
                    //this should use action dispatch to update the state
                    var suffix = "";
                    if (that.props.meta){
                      suffix = "<br/> " + that.props.meta[i];
                    }

                    var t = datetimeFormat(new Date(xdata[i] * 1000));
                    JQuery(React.findDOMNode(that))
                                .find('.bar-value')
                                .html(t + " → " + yFormatter(d) + suffix);
                  })
                  .on("mouseout", (d, i) => {
                    //this should use action dispatch to update the state
                    JQuery(React.findDOMNode(that)).find('.bar-value').text("");
                  });

    svg.append("g")
       .attr("class", "x axis")
       .attr("transform", "translate("+ margin.left+"," + height + ")")
       .call(xAxis)
       .selectAll("text")
       .attr("dx", "-.8em")
       .attr("dy", ".15em")
       .style("text-anchor", "end")
       .attr("transform", "rotate(-65)" );

    svg.append("g")
       .attr("class", "y axis")
       .attr("transform", "translate("+ margin.left + ", 0)")
       .call(yAxis);

    var subtitle = "";
    if (this.props.subtitle){
      subtitle = (<h6>{this.props.subtitle}</h6>);

    }
    return (
      <div className="chart">
        <div className="bar-value text-right">
        </div>
        <div className="row">
          <div className="col-sm-8">
            <h5>{this.props.title}</h5>
          </div>
        </div>
        <div className="row">
          <div className="col-sm-12">
            {subtitle}
          </div>
        </div>
        {node.toReact()}
      </div>
    );

  }
}

export default Chart;
